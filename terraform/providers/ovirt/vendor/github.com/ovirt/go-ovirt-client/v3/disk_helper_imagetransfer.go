package ovirtclient

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// correlationIDRand is a random generator for selecting letters to put in the correlation ID. This does not need
// to be cryptographically strong as it is short-lived.
var correlationIDRand = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec,gochecknoglobals

// generateCorrelationID generates a random ID usable for correlation.
func generateCorrelationID(prefix string) string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[correlationIDRand.Intn(len(letters))]
	}
	return fmt.Sprintf("%s%s", prefix, string(b))
}

// newImageTransfer creates a new image transfer for both uploads and downloads of images. It must be passed the
// following parameters:
//
//   - cli is the oVirt SDK client.
//   - logger is a logger from the go-ovirt-client-logger library
//   - diskID is the ID of the disk that is being transferred to/from.
//   - correlationID is an optional unique ID that can be used to check if the job completed. If no correlation ID is
//     passed, this function generates a random one.
//   - retries is a list of retry strategies to use for each API call.
//   - direction is the direction of transfer. See ovirtsdk.ImageTransferDirection.
//   - format is the disk format being uploaded, or the disk format requested. The oVirt Engine will automatically convert
//     images to the requested format.
//   - updateDisk is a function that will be called whenever the disk object is updated.
func newImageTransfer(
	cli *oVirtClient,
	logger Logger,
	diskID DiskID,
	correlationID string,
	retries []RetryStrategy,
	direction ovirtsdk4.ImageTransferDirection,
	format ovirtsdk4.DiskFormat,
	updateDisk func(disk Disk),
) imageTransfer {
	if correlationID == "" {
		correlationID = generateCorrelationID(fmt.Sprintf("image_%s_", direction))
	}

	return &imageTransferImpl{
		retries:         retries,
		diskID:          diskID,
		cli:             cli,
		logger:          logger,
		correlationID:   correlationID,
		conn:            cli.conn,
		transfer:        nil,
		transferService: nil,
		httpClient:      cli.httpClient,
		direction:       direction,
		format:          format,
		updateDisk:      updateDisk,
	}
}

// imageTransfer is an internal helper to facilitate image transfers from/to the oVirt Engine. It should not be reused
// for multiple transfers.
type imageTransfer interface {
	// initialize sets up the image transfer in the specified direction. If successful, it returns
	// the URL the image needs to be transferred to/from using a HTTP request. If not, it aborts the
	// transfer and returns an error. In any case, the calling party MUST call finalize to correctly
	// finalize the image transfer.
	initialize() (transferURL string, err error)
	// finalize cleans up the image transfer. It must be called regardless if an error happened or not, and the error
	// must be passed to it so it can determine how to best clean up the image transfer.
	//
	// If the finalize function is not called the disk may potentially stay in locked status indefinitely.
	finalize(err error) error

	// checkStatusCode checks an ImageIO status code for correctness and returns an error if it is
	// not correct.
	checkStatusCode(statusCode int) error
}

// imageTransferImpl is the implementation of the imageTransfer interface.
type imageTransferImpl struct {
	// retries is the list of retry strategies to use for calls in this transfer.
	retries []RetryStrategy
	// diskID is the ID of the disk used for this transfer.
	diskID DiskID
	// cli is the calling client library.
	cli *oVirtClient
	// logger is the go-ovirt-client-log logger
	logger Logger
	// correlationID is a unique ID that can be used to track jobs in the oVirt Engine.
	correlationID string
	// conn is the underlying SDK connection.
	conn *ovirtsdk4.Connection
	// httpClient is the configured HTTP client for calling the engine.
	httpClient http.Client
	// direction indicates the direction of transfer.
	direction ovirtsdk4.ImageTransferDirection
	// format is the image format for the transfer. For downloads, the engine will convert the image to
	// this format. For uploads, the engine will expect the image to be send in this format.
	format ovirtsdk4.DiskFormat
	// updateDisk is a callback that allows a calling party to get notified when the underlying disk object
	// changes.
	updateDisk func(disk Disk)
	// transfer is the created image transfer. It is set after the createImageTransfer function
	// is called.
	transfer *ovirtsdk4.ImageTransfer
	// transferService is the service associated with transfer. it is set after the createImageTransfer
	// function is called.
	transferService *ovirtsdk4.ImageTransferService
	// transferURL is the URL that is found for the transfer. It is set after findTransferURL is called.
	transferURL string
}

// checkStatusCode takes a HTTP status code from the ImageIO endpoint and verifies it.
func (i *imageTransferImpl) checkStatusCode(statusCode int) error {
	if statusCode < 300 {
		return nil
	}
	switch {
	case statusCode < 399:
		return newError(
			ENotAnOVirtEngine,
			"received redirect response for image %s",
			i.direction,
		)
	case statusCode < 499:
		if statusCode == 401 {
			return newError(
				EAccessDenied,
				"received unauthorized (401) status code for image %s",
				i.direction,
			)
		}
		return newError(
			EPermanentHTTPError,
			"unexpected client status code (%d) received for image %s",
			statusCode,
			i.direction,
		)
	default:
		return newError(
			EPermanentHTTPError,
			"unexpected server error status code %d while attempting to %s image",
			statusCode,
			i.direction,
		)
	}
}

// initialize sets up the image transfer in the specified direction. If successful, it returns
// the URL the image needs to be transferred to/from using a HTTP request. If not, it aborts the
// transfer and returns an error. In any case, the calling party MUST call finalize to correctly
// finalize the image transfer.
func (i *imageTransferImpl) initialize() (transferURL string, err error) {
	steps := []func() error{
		i.waitForTransferOk,
		i.createImageTransfer,
		i.waitForImageTransferReady,
		i.findTransferURL,
	}

	for _, step := range steps {
		if err := step(); err != nil {
			i.abortTransfer()
			return "", err
		}
	}
	return i.transferURL, nil
}

// finalize finalizes or aborts the image transfer, depending on if an error happened. The calling
// party must pass any error that happened so that the finalize function can make the correct decision.
func (i *imageTransferImpl) finalize(err error) error {
	if err != nil {
		i.abortTransfer()
		return err
	}
	steps := []func() error{
		i.finalizeTransfer,
		i.waitForTransferFinalize,
		i.waitForTransferOk,
	}
	for _, step := range steps {
		if err := step(); err != nil {
			i.abortTransfer()
			return err
		}
	}
	return nil
}

// waitForTransferOk waits for a disk to be in the OK status, then additionally queries the job that was in progress with
// the correlation ID. This is necessary because the disk returns OK status before the job has actually finished,
// resulting in a "disk locked" error on subsequent operations. It uses checkDiskOk as an underlying function.
//
// This function also calls the updateDisk hook to update the disk on the calling side.
func (i *imageTransferImpl) waitForTransferOk() (err error) {
	disk, err := i.cli.WaitForDiskOK(i.diskID, i.retries...)

	if err != nil {
		return err
	}

	if err := i.cli.waitForJobFinished(i.correlationID, i.retries); err != nil {
		return err
	}

	i.updateDisk(disk)
	return nil
}

// buildImageTransferRequest creates an SDK image transfer request and the associated service.
func (i *imageTransferImpl) buildImageTransferRequest() (
	*ovirtsdk4.ImageTransfersServiceAddRequest,
	*ovirtsdk4.ImageTransfersService,
) {
	imageTransfersService := i.conn.SystemService().ImageTransfersService()
	image := ovirtsdk4.NewImageBuilder().Id(string(i.diskID)).MustBuild()
	transfer := ovirtsdk4.
		NewImageTransferBuilder().
		Image(image).
		Direction(i.direction).
		Format(i.format).
		MustBuild()
	transferReq := imageTransfersService.
		Add().
		ImageTransfer(transfer).
		Query("correlation_id", i.correlationID)
	return transferReq, imageTransfersService
}

// createImageTransfer repeatedly tries to create an image transfer until it succeeds or it runs out of retries.
// This function will set the i.transfer and i.transferService variables with the created image transfer and
// the associated service.
func (i *imageTransferImpl) createImageTransfer() (err error) {
	return retry(
		fmt.Sprintf("starting image transfer for disk %s", i.diskID),
		i.logger,
		i.retries,
		i.attemptCreateImageTransfer,
	)
}

// attemptCreateImageTransfer attempts to create an image transfer with the oVirt Engine API and returns an error if it
// fails. createImageTransfer can be used to repeatedly try this function.
func (i *imageTransferImpl) attemptCreateImageTransfer() error {
	transferReq, imageTransfersService := i.buildImageTransferRequest()

	transferRes, e := transferReq.Send()
	if e != nil {
		return e
	}
	var ok bool
	i.transfer, ok = transferRes.ImageTransfer()
	if !ok {
		return newError(
			EFieldMissing,
			"missing image transfer as a response to image transfer create request",
		)
	}
	transferID, ok := i.transfer.Id()
	if !ok {
		return newError(
			EFieldMissing,
			"missing image transfer ID in response to image transfer create request",
		)
	}
	i.transferService = imageTransfersService.ImageTransferService(transferID)
	return nil
}

// waitForImageTransferReady repeatedly calls checkImageTransferReady until it returns successfully or the retries are
// exhausted.
//
// This function is internal to imageTransferImpl, do not call externally.
func (i *imageTransferImpl) waitForImageTransferReady() (err error) {
	return retry(
		fmt.Sprintf(
			"waiting for image transfer to become ready for disk ID %s",
			i.diskID,
		),
		i.logger,
		i.retries,
		i.checkImageTransferReady,
	)
}

// checkImageTransferReady retrieves the image transfer once and checks if it is in the transferring phase.
// waitForImageTransferReady can be used to call this function repeatedly.
func (i *imageTransferImpl) checkImageTransferReady() error {
	req, err := i.transferService.Get().Send()
	if err != nil {
		return err
	}
	transfer, ok := req.ImageTransfer()
	if !ok {
		return newError(
			EFieldMissing,
			"fetching image transfer did not return an image transfer",
		)
	}
	phase, ok := transfer.Phase()
	if !ok {
		return newError(
			EFieldMissing,
			"fetching image transfer did not contain a phase",
		)
	}
	switch phase {
	case ovirtsdk4.IMAGETRANSFERPHASE_INITIALIZING:
		return newError(
			EPending,
			"image transfer is in phase %s instead of transferring",
			phase,
		)
	case ovirtsdk4.IMAGETRANSFERPHASE_TRANSFERRING:
		return nil
	default:
		return newError(
			EUnexpectedImageTransferPhase,
			"image transfer is in phase %s instead of %s",
			phase,
			ovirtsdk4.IMAGETRANSFERPHASE_TRANSFERRING,
		)
	}
}

// finalizeTransfer attempts to finalize an image transfer. It is part of the finalize function. After this function
// finalize still waits for the disk to be OK. This function calls attemptFinalizeTransfer repeatedly until it succeeds
// or the retries are exhausted.
func (i *imageTransferImpl) finalizeTransfer() error {
	return retry(
		fmt.Sprintf("finalizing image for disk %s", i.diskID),
		i.logger,
		i.retries,
		i.attemptFinalizeTransfer,
	)
}

// attemptFinalizeTransfer calls the oVirt Engine API a single time attempting to finalize a transfer.
func (i *imageTransferImpl) attemptFinalizeTransfer() error {
	finalizeRequest := i.transferService.Finalize()
	finalizeRequest.Query("correlation_id", i.correlationID)
	_, err := finalizeRequest.Send()
	return err
}

// waitForTransferFinalize waits for a transfer to reach a final state.
func (i *imageTransferImpl) waitForTransferFinalize() error {
	return retry(
		fmt.Sprintf("waiting for finalizing image transfer for disk %s", i.diskID),
		i.logger,
		i.retries,
		i.attemptWaitForTransferFinalize,
	)
}

// attemptWaitForTransferFinalize runs a single request checking for the transfer to finalize.
func (i *imageTransferImpl) attemptWaitForTransferFinalize() error {
	return i.checkImageTransferPhase(
		ovirtsdk4.IMAGETRANSFERPHASE_FINISHED_SUCCESS,
		[]ovirtsdk4.ImageTransferPhase{
			ovirtsdk4.IMAGETRANSFERPHASE_FINISHED_FAILURE,
			ovirtsdk4.IMAGETRANSFERPHASE_PAUSED_SYSTEM,
		},
	)
}

func (i *imageTransferImpl) waitForTransferAbort() error {
	return retry(
		fmt.Sprintf("waiting for aborting image transfer for disk %s", i.diskID),
		i.logger,
		i.retries,
		i.attemptWaitForTransferAbort,
	)
}

// attemptWaitForTransferAbort runs a single request checking for the transfer to abort.
func (i *imageTransferImpl) attemptWaitForTransferAbort() error {
	return i.checkImageTransferPhase(
		ovirtsdk4.IMAGETRANSFERPHASE_FINISHED_FAILURE,
		[]ovirtsdk4.ImageTransferPhase{
			ovirtsdk4.IMAGETRANSFERPHASE_FINISHED_SUCCESS,
			ovirtsdk4.IMAGETRANSFERPHASE_PAUSED_SYSTEM,
		},
	)
}

func (i *imageTransferImpl) checkImageTransferPhase(
	waitForPhase ovirtsdk4.ImageTransferPhase,
	disallowedPhases []ovirtsdk4.ImageTransferPhase,
) error {
	var notFoundError *ovirtsdk4.NotFoundError
	transferResponse, err := i.transferService.Get().Send()
	if err != nil {
		if errors.As(err, &notFoundError) {
			// The image transfer disappeared, which happens on oVirt <4.4.7. The calling
			// party must now wait for the disk to be OK, nothing left for us to do.
			return nil
		}
		return err
	}
	transfer, ok := transferResponse.ImageTransfer()
	if !ok {
		// Image transfer has disappeared, see comment above.
		return nil
	}
	if transfer.MustPhase() == waitForPhase {
		return nil
	}
	for _, phase := range disallowedPhases {
		if transfer.MustPhase() == phase {
			return newError(
				EUnexpectedImageTransferPhase,
				"Unexpected image transfer phase %s",
				transfer.MustPhase(),
			)
		}
	}
	return newError(EPending, "Transfer is in phase %s", transfer.MustPhase())
}

// findTransferURL sends HTTP OPTIONS requests to potential transfer URLs via verifyTransferURL to determine if a
// transfer URL can be used or not. This method sets the i.transferURL variable.
func (i *imageTransferImpl) findTransferURL() (err error) {
	i.logger.Debugf(
		"Attempting to determine image transfer URL for disk %s...",
		i.diskID,
	)
	var tryURLs []string
	if transferURL, ok := i.transfer.TransferUrl(); ok && transferURL != "" {
		tryURLs = append(tryURLs, transferURL)
	}
	if proxyURL, ok := i.transfer.ProxyUrl(); ok && proxyURL != "" {
		tryURLs = append(tryURLs, proxyURL)
	}

	if len(tryURLs) == 0 {
		i.logger.Errorf(
			"Bug: neither a transfer URL nor a proxy URL was returned from the oVirt Engine. (%v)",
			i.transfer,
		)
		return newError(EBug, "neither a transfer URL nor a proxy URL was returned from the oVirt Engine")
	}

	var lastError error
	for _, transferURL := range tryURLs {
		lastError = i.verifyTransferURL(transferURL)
		if lastError == nil {
			i.transferURL = transferURL
			return nil
		}
	}
	if lastError != nil {
		return wrap(
			lastError,
			EConnection,
			"failed to find a valid transfer URL; check your network connectivity to the oVirt Engine ImageIO port",
		)
	}
	return nil
}

// verifyTransferURL takes a transfer URL from findTransferURL and repeatedly sends the OPTIONS request via
// optionsRequest to figure out if the URL can be used. It tries a maximum of 3 times.
func (i *imageTransferImpl) verifyTransferURL(transferURL string) error {
	parsedTransferURL, err := url.Parse(transferURL)
	if err != nil {
		return wrap(err, EUnidentified, "failed to parse transfer URL %s", transferURL)
	}

	return retry(
		fmt.Sprintf("sending OPTIONS request to %s", transferURL),
		i.logger,
		append(i.retries, MaxTries(3)),
		func() error {
			return i.optionsRequest(parsedTransferURL)
		},
	)
}

// optionsRequest sends an individual options request to the specified URL to figure out if the URL can be used for
// an image transfer.
func (i *imageTransferImpl) optionsRequest(parsedTransferURL *url.URL) error {
	optionsReq, e := http.NewRequest(http.MethodOptions, parsedTransferURL.String(), strings.NewReader(""))
	if e != nil {
		return wrap(e, EBug, "failed to create OPTIONS request to %s", parsedTransferURL.String())
	}
	res, e := i.httpClient.Do(optionsReq)
	if e != nil {
		return wrap(e, EConnection, "HTTP request to %s failed", parsedTransferURL.String())
	}
	defer func() {
		_ = res.Body.Close()
	}()
	statusCode := res.StatusCode
	switch {
	case statusCode < 199:
		return newError(
			EConnection,
			"HTTP connection error while calling %s",
			parsedTransferURL.String(),
		)
	case statusCode < 399:
		return nil
	case statusCode < 499:
		return newError(
			EPermanentHTTPError,
			"HTTP 4xx status code returned from URL %s (%d)",
			parsedTransferURL.String(),
			res.StatusCode,
		)
	default:
		return newError(
			EConnection,
			"non-200 status code returned from URL %s (%d)",
			parsedTransferURL.String(),
			res.StatusCode,
		)
	}
}

// abortTransfer cancels an image transfer with the oVirt Engine API. It calls the abort repeatedly until it succeeds or
// the retries are exhausted.
func (i *imageTransferImpl) abortTransfer() {
	if i.transfer != nil {
		errorHappened := false
		if err := retry(
			fmt.Sprintf("canceling transfer for disk %s", i.diskID),
			i.logger,
			i.retries,
			i.attemptAbortTransfer,
		); err != nil {
			// We can't really do anything as we are already in a failure state, log the error.
			i.logger.Warningf(
				"failed to cancel transfer for disk %s, may not be able to remove disk",
				i.diskID,
			)
			errorHappened = true
		}
		if err := i.waitForTransferAbort(); err != nil {
			i.logger.Warningf(
				"failed to wait for disk %s to return to OK state after aborting transfer",
				i.diskID,
			)
			errorHappened = true
		}
		if err := i.waitForTransferOk(); err != nil && !HasErrorCode(err, ENotFound) {
			// We can't really do anything as we are already in a failure state, log the error.
			// The ENotFound is expected as the disk may be removed if a transfer is aborted.
			i.logger.Warningf(
				"failed to wait for disk %s to return to OK state after aborting transfer",
				i.diskID,
			)
			errorHappened = true
		}
		if !errorHappened {
			i.transfer = nil
		}
	}
}

// attemptAbortTransfer attempts to cancel an image transfer with the oVirt Engine API.
func (i *imageTransferImpl) attemptAbortTransfer() error {
	_, err := i.transferService.Cancel().Send()
	return err
}
