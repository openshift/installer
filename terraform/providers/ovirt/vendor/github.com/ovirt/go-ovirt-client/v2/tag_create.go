package ovirtclient

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CreateTag(name string, params CreateTagParams, retries ...RetryStrategy) (result Tag, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	if params == nil {
		params = NewCreateTagParams()
	}

	err = retry(
		"creating tag",
		o.logger,
		retries,
		func() error {
			tagBuilder := ovirtsdk.NewTagBuilder().Name(name)
			if description := params.Description(); description != nil {
				tagBuilder.Description(*description)
			}
			response, e := o.conn.SystemService().TagsService().Add().Tag(tagBuilder.MustBuild()).Send()
			if e != nil {
				return e
			}

			tag, ok := response.Tag()
			if !ok {
				return newError(EFieldMissing, "missing Tag in response")
			}

			result, err = convertSDKTag(tag, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert Tag",
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) CreateTag(name string, params CreateTagParams, _ ...RetryStrategy) (result Tag, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	id := TagID(m.GenerateUUID())
	if params == nil {
		params = NewCreateTagParams()
	}
	tag := &tag{
		client:      m,
		id:          id,
		name:        name,
		description: params.Description(),
	}
	m.tags[id] = tag

	result = tag
	return
}
