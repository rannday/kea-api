package client

// BuildReport fetches the build-report for a single service.
func BuildReport(c *Client, service Service) (string, error) {
	return CallAndExtractText(c, "build-report", service)
}

// BuildReportMulti fetches the build-report for multiple services.
func BuildReportMulti(c *Client, services ...Service) ([]string, error) {
	responses, err := CallCommand(c, "build-report", services...)
	if err != nil {
		return nil, err
	}

	reports := make([]string, 0, len(responses))
	for _, res := range responses {
		reports = append(reports, res.Text)
	}
	return reports, nil
}

// ConfigGet fetches the config for a service and decodes it into T.
func ConfigGet[T any](c *Client, service Service) (T, error) {
	return DecodeFirst[T](c, "config-get", service)
}

// ConfigGetMulti fetches and decodes the configuration of multiple services into []T.
func ConfigGetMulti[T any](c *Client, services ...Service) ([]T, error) {
	return CallAndDecode[T](c, "config-get", services...)
}

// ListCommands fetches the list of supported commands for a service.
func ListCommands(c *Client, service Service) ([]string, error) {
	return DecodeFirst[[]string](c, "list-commands", service)
}

// ListCommandsMulti fetches the list of supported commands for multiple services.
func ListCommandsMulti(c *Client, services ...Service) ([][]string, error) {
	return CallAndDecode[[]string](c, "list-commands", services...)
}

// StatusGet fetches status information for a service and decodes into T.
func StatusGet[T any](c *Client, service Service) (T, error) {
	return DecodeFirst[T](c, "status-get", service)
}

// StatusGetMulti fetches status information for multiple services and decodes into []T.
func StatusGetMulti[T any](c *Client, services ...Service) ([]T, error) {
	return CallAndDecode[T](c, "status-get", services...)
}

// VersionGet fetches version info for a service, returning both the full text and the decoded T.
func VersionGet[T any](c *Client, service Service) (string, T, error) {
	return DecodeFirstWithText[T](c, "version-get", service)
}

// VersionGetMulti fetches version info for multiple services, returning the raw text for each.
func VersionGetMulti(c *Client, services ...Service) ([]string, error) {
	responses, err := CallCommand(c, "version-get", services...)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, res := range responses {
		versions = append(versions, res.Text)
	}
	return versions, nil
}
