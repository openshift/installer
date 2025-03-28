package gcp

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/googleapis/gax-go/v2/apierror"
	googleoauth "golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/googleapi"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset/installconfig/gcp/mock"
	"github.com/openshift/installer/pkg/types/gcp"
)

var (
	testTags = []gcp.UserTag{
		{ParentID: "openshift", Key: "key1", Value: "value1"}, {ParentID: "openshift", Key: "key2", Value: "value2"},
		{ParentID: "openshift", Key: "key3", Value: "value3"}, {ParentID: "openshift", Key: "key4", Value: "value4"},
		{ParentID: "openshift", Key: "key5", Value: "value5"}, {ParentID: "openshift", Key: "key6", Value: "value6"},
		{ParentID: "openshift", Key: "key7", Value: "value7"}, {ParentID: "openshift", Key: "key8", Value: "value8"},
		{ParentID: "openshift", Key: "key9", Value: "value9"}, {ParentID: "openshift", Key: "key10", Value: "value10"},
		{ParentID: "openshift", Key: "key11", Value: "value11"}, {ParentID: "openshift", Key: "key12", Value: "value12"},
		{ParentID: "openshift", Key: "key13", Value: "value13"}, {ParentID: "openshift", Key: "key14", Value: "value14"},
		{ParentID: "openshift", Key: "key15", Value: "value15"}, {ParentID: "openshift", Key: "key16", Value: "value16"},
		{ParentID: "openshift", Key: "key17", Value: "value17"}, {ParentID: "openshift", Key: "key18", Value: "value18"},
		{ParentID: "openshift", Key: "key19", Value: "value19"}, {ParentID: "openshift", Key: "key20", Value: "value20"},
		{ParentID: "openshift", Key: "key21", Value: "value21"}, {ParentID: "openshift", Key: "key22", Value: "value22"},
		{ParentID: "openshift", Key: "key23", Value: "value23"}, {ParentID: "openshift", Key: "key24", Value: "value24"},
		{ParentID: "openshift", Key: "key25", Value: "value25"}, {ParentID: "openshift", Key: "key26", Value: "value26"},
		{ParentID: "openshift", Key: "key27", Value: "value27"}, {ParentID: "openshift", Key: "key28", Value: "value28"},
		{ParentID: "openshift", Key: "key29", Value: "value29"}, {ParentID: "openshift", Key: "key30", Value: "value30"},
		{ParentID: "openshift", Key: "key31", Value: "value31"}, {ParentID: "openshift", Key: "key32", Value: "value32"},
		{ParentID: "openshift", Key: "key33", Value: "value33"}, {ParentID: "openshift", Key: "key34", Value: "value34"},
		{ParentID: "openshift", Key: "key35", Value: "value35"}, {ParentID: "openshift", Key: "key36", Value: "value36"},
		{ParentID: "openshift", Key: "key37", Value: "value37"}, {ParentID: "openshift", Key: "key38", Value: "value38"},
		{ParentID: "openshift", Key: "key39", Value: "value39"}, {ParentID: "openshift", Key: "key40", Value: "value40"},
		{ParentID: "openshift", Key: "key41", Value: "value41"}, {ParentID: "openshift", Key: "key42", Value: "value42"},
		{ParentID: "openshift", Key: "key43", Value: "value43"}, {ParentID: "openshift", Key: "key44", Value: "value44"},
		{ParentID: "openshift", Key: "key45", Value: "value45"}, {ParentID: "openshift", Key: "key46", Value: "value46"},
		{ParentID: "openshift", Key: "key47", Value: "value47"}, {ParentID: "openshift", Key: "key48", Value: "value48"},
		{ParentID: "openshift", Key: "key49", Value: "value49"}, {ParentID: "openshift", Key: "key50", Value: "value50"},
		{ParentID: "openshift", Key: "key51", Value: "value51"}, {ParentID: "openshift", Key: "key52", Value: "value52"},
	}

	testProcessedTags = map[string]struct {
		key   string
		value string
	}{
		"openshift/key1/value1":   {key: "tagKeys/0005159847", value: "tagValues/0030662057"},
		"openshift/key2/value2":   {key: "tagKeys/0018550390", value: "tagValues/0073384808"},
		"openshift/key3/value3":   {key: "tagKeys/0051315427", value: "tagValues/0051662048"},
		"openshift/key4/value4":   {key: "tagKeys/0081126737", value: "tagValues/0012526301"},
		"openshift/key5/value5":   {key: "tagKeys/0009228350", value: "tagValues/0043054391"},
		"openshift/key6/value6":   {key: "tagKeys/0051626174", value: "tagValues/0027708252"},
		"openshift/key7/value7":   {key: "tagKeys/0030414034", value: "tagValues/0020489890"},
		"openshift/key8/value8":   {key: "tagKeys/0050469265", value: "tagValues/0002904904"},
		"openshift/key9/value9":   {key: "tagKeys/0076293608", value: "tagValues/0055629312"},
		"openshift/key10/value10": {key: "tagKeys/0023033386", value: "tagValues/0088842980"},
		"openshift/key11/value11": {key: "tagKeys/0069899625", value: "tagValues/0047417661"},
		"openshift/key12/value12": {key: "tagKeys/0034605069", value: "tagValues/0050384905"},
		"openshift/key13/value13": {key: "tagKeys/0028357547", value: "tagValues/0052268968"},
		"openshift/key14/value14": {key: "tagKeys/0099944474", value: "tagValues/0059052883"},
		"openshift/key15/value15": {key: "tagKeys/0050205103", value: "tagValues/0036307885"},
		"openshift/key16/value16": {key: "tagKeys/0079771629", value: "tagValues/0065673174"},
		"openshift/key17/value17": {key: "tagKeys/0060225722", value: "tagValues/0081145498"},
		"openshift/key18/value18": {key: "tagKeys/0016496476", value: "tagValues/0046494994"},
		"openshift/key19/value19": {key: "tagKeys/0093247819", value: "tagValues/0041540373"},
		"openshift/key20/value20": {key: "tagKeys/0080859513", value: "tagValues/0016693395"},
		"openshift/key21/value21": {key: "tagKeys/0018537779", value: "tagValues/0003454649"},
		"openshift/key22/value22": {key: "tagKeys/0071724280", value: "tagValues/0047292544"},
		"openshift/key23/value23": {key: "tagKeys/0045095645", value: "tagValues/0089378558"},
		"openshift/key24/value24": {key: "tagKeys/0044575217", value: "tagValues/0022754275"},
		"openshift/key25/value25": {key: "tagKeys/0056084774", value: "tagValues/0040808197"},
		"openshift/key26/value26": {key: "tagKeys/0086508506", value: "tagValues/0091979350"},
		"openshift/key27/value27": {key: "tagKeys/0085330359", value: "tagValues/0051833259"},
		"openshift/key28/value28": {key: "tagKeys/0094744916", value: "tagValues/0011642000"},
		"openshift/key29/value29": {key: "tagKeys/0014270555", value: "tagValues/0072404680"},
		"openshift/key30/value30": {key: "tagKeys/0079085850", value: "tagValues/0007793185"},
		"openshift/key31/value31": {key: "tagKeys/0031484153", value: "tagValues/0050294705"},
		"openshift/key32/value32": {key: "tagKeys/0045311563", value: "tagValues/0029329808"},
		"openshift/key33/value33": {key: "tagKeys/0080836115", value: "tagValues/0003514535"},
		"openshift/key34/value34": {key: "tagKeys/0072216154", value: "tagValues/0060486146"},
		"openshift/key35/value35": {key: "tagKeys/0025032284", value: "tagValues/0038979234"},
		"openshift/key36/value36": {key: "tagKeys/0057998529", value: "tagValues/0067716498"},
		"openshift/key37/value37": {key: "tagKeys/0086808493", value: "tagValues/0060060909"},
		"openshift/key38/value38": {key: "tagKeys/0029402635", value: "tagValues/0060648494"},
		"openshift/key39/value39": {key: "tagKeys/0034805062", value: "tagValues/0066064364"},
		"openshift/key40/value40": {key: "tagKeys/0052208445", value: "tagValues/0098376015"},
		"openshift/key41/value41": {key: "tagKeys/0078738427", value: "tagValues/0046807958"},
		"openshift/key42/value42": {key: "tagKeys/0075526710", value: "tagValues/0092960786"},
		"openshift/key43/value43": {key: "tagKeys/0071238172", value: "tagValues/0014715775"},
		"openshift/key44/value44": {key: "tagKeys/0098580341", value: "tagValues/0052744277"},
		"openshift/key45/value45": {key: "tagKeys/0046775969", value: "tagValues/0095864916"},
		"openshift/key46/value46": {key: "tagKeys/0041543652", value: "tagValues/0054302656"},
		"openshift/key47/value47": {key: "tagKeys/0043236402", value: "tagValues/0087723355"},
		"openshift/key48/value48": {key: "tagKeys/0098289690", value: "tagValues/0018491792"},
		"openshift/key49/value49": {key: "tagKeys/0022125412", value: "tagValues/0063407483"},
		"openshift/key50/value50": {key: "tagKeys/0004322107", value: "tagValues/0083960428"},
		"openshift/key51/value51": {key: "tagKeys/0070391549", value: "tagValues/0027514043"},
		"openshift/key51/value52": {key: "tagKeys/0070391542", value: "tagValues/0227514043"},
	}
)

func generateExpectedProcessedTags(tags []gcp.UserTag) map[string]string {
	retTags := make(map[string]string, len(tags))
	for _, tag := range tags {
		t := testProcessedTags[fmt.Sprintf("%s/%s/%s", tag.ParentID, tag.Key, tag.Value)]
		retTags[t.key] = t.value
	}
	return retTags
}

func getTestGetNamespacedTagValueForbiddenError(tagValueNamespacedName string) error {
	apiErr, _ := apierror.FromError(fmt.Errorf("%w", &googleapi.Error{
		Code:    http.StatusForbidden,
		Message: fmt.Sprintf("Permission denied on resource '%s' (or it may not exist).", tagValueNamespacedName),
	}))
	return apiErr
}

func getTestGetNamespacedTagValueInternalError(tagValueNamespacedName string) error {
	apiErr, _ := apierror.FromError(fmt.Errorf("%w", &googleapi.Error{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf("Internal error while fetching '%s'", tagValueNamespacedName),
	}))
	return apiErr
}

func getTestNamespacedTagValueResp(tagValueNamespacedName string) *cloudresourcemanager.TagValue {
	tag := testProcessedTags[tagValueNamespacedName]
	return &cloudresourcemanager.TagValue{
		Name:           tag.value,
		Parent:         tag.key,
		NamespacedName: tagValueNamespacedName,
	}
}

func resetProcessedTags() {
	processedTags.Lock()
	defer processedTags.Unlock()

	processedTags.processed = false
	processedTags.tags = map[string]string{}
}

func TestGetUserTags(t *testing.T) {
	testProject := "test-project"
	// serializer is for serializing the access to processedTags, since
	// each scenario expects different set of tags.
	serializer := &sync.Mutex{}

	testCases := []struct {
		name          string
		projectID     string
		userTags      []gcp.UserTag
		processedTags map[string]string
		expectedError string
	}{
		{
			name:          "user tags is empty",
			userTags:      []gcp.UserTag{},
			processedTags: map[string]string{},
		},
		{
			name:          "more than max allowed tags configured",
			userTags:      testTags,
			processedTags: nil,
			expectedError: "more than 50 user tags is not allowed, configured count: 52",
		},
		{
			name:          "user tags processing passed",
			userTags:      testTags[2:50],
			processedTags: generateExpectedProcessedTags(testTags[2:50]),
		},
		{
			name:          "fetching project tags fails",
			projectID:     "fail-project",
			userTags:      testTags[:5],
			processedTags: nil,
			expectedError: "failed to fetch tags attached to fail-project: The caller does not have permission",
		},
		{
			name:          "found duplicate tags",
			projectID:     "dup-project",
			userTags:      testTags[5:10],
			processedTags: nil,
			expectedError: `found duplicate tags, [openshift/key6/value6 openshift/key7/value7 openshift/key8/value8 openshift/key9/value9 openshift/key10/value10] tags already exist on dup-project project resource`,
		},
		{
			name:          "non-existent tags configured",
			userTags:      testTags[40:51],
			processedTags: nil,
			expectedError: `does not have permission to access [openshift/key51/value51] tag(s) or does not exist`,
		},
		{
			name:          "failed to fetch tag",
			userTags:      testTags[40:],
			processedTags: nil,
			expectedError: `failed to fetch user-defined tag openshift/key52/value52: googleapi: Error 500: Internal error while fetching 'openshift/key52/value52'`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gcpClient := mock.NewMockAPI(mockCtrl)
	tagMgr := NewTagManager(gcpClient)

	// Return fake credentials.
	gcpClient.EXPECT().GetCredentials().Return(&googleoauth.Credentials{JSON: []byte(fakeCreds)}).AnyTimes()

	// Return ["openshift/key1/value1", "openshift/key2/value2"] tags for "test-project" project.
	gcpClient.EXPECT().GetProjectTags(gomock.Any(), testProject).Return(sets.Insert(sets.New[string](), "openshift/key1/value1", "openshift/key2/value2"), nil).AnyTimes()

	// Return error for "fail-project" project.
	gcpClient.EXPECT().GetProjectTags(gomock.Any(), "fail-project").
		Return(
			nil,
			fmt.Errorf("failed to fetch tags attached to fail-project: The caller does not have permission")).
		AnyTimes()

	// Return ["openshift/key6/value6", "openshift/key7/value7", "openshift/key8/value8",
	// "openshift/key9/value9", "openshift/key10/value10"] tags for "dup-project" project.
	gcpClient.EXPECT().GetProjectTags(gomock.Any(), "dup-project").
		Return(
			sets.Insert(sets.New[string](),
				"openshift/key6/value6",
				"openshift/key7/value7",
				"openshift/key8/value8",
				"openshift/key9/value9",
				"openshift/key10/value10",
			), nil).AnyTimes()

	// Return Forbidden error "openshift/key51/value51" tag.
	gcpClient.EXPECT().GetNamespacedTagValue(gomock.Any(), "openshift/key51/value51").Return(nil, getTestGetNamespacedTagValueForbiddenError("openshift/key51/value51")).AnyTimes()

	// Return Internal error "openshift/key52/value52" tag.
	gcpClient.EXPECT().GetNamespacedTagValue(gomock.Any(), "openshift/key52/value52").Return(nil, getTestGetNamespacedTagValueInternalError("openshift/key52/value52")).AnyTimes()

	// Return details for the requested tag.
	gcpClient.EXPECT().GetNamespacedTagValue(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, tagNamespacedName string) (*cloudresourcemanager.TagValue, error) {
			return getTestNamespacedTagValueResp(tagNamespacedName), nil
		}).AnyTimes()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serializer.Lock()
			defer serializer.Unlock()
			resetProcessedTags()

			projectID := testProject
			if tc.projectID != "" {
				projectID = tc.projectID
			}
			tags, err := tagMgr.GetUserTags(context.Background(), projectID, tc.userTags)

			if (err != nil || tc.expectedError != "") && err.Error() != tc.expectedError {
				t.Errorf("GetUserTags(): error: got: %v, want: %v", err, tc.expectedError)
			}

			if !reflect.DeepEqual(tags, tc.processedTags) {
				t.Errorf("GetUserTags(): tags: got: %v, want: %v", tags, tc.processedTags)
			}
		})
	}
}
