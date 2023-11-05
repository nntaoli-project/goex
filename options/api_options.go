package options

type ApiOptions struct {
	Key        string
	Secret     string
	Passphrase string
	ClientId   string
	Simulated  bool
}

type ApiOption func(options *ApiOptions)

func WithApiKey(key string) ApiOption {
	return func(options *ApiOptions) {
		options.Key = key
	}
}

func WithApiSecretKey(secret string) ApiOption {
	return func(options *ApiOptions) {
		options.Secret = secret
	}
}

func WithPassphrase(passphrase string) ApiOption {
	return func(options *ApiOptions) {
		options.Passphrase = passphrase
	}
}

func WithClientId(clientId string) ApiOption {
	return func(options *ApiOptions) {
		options.ClientId = clientId
	}
}

func WithSimulated(isSimulated ...bool) ApiOption {
	return func(options *ApiOptions) {
		if len(isSimulated) <= 0 {
			options.Simulated = true
		} else {
			options.Simulated = isSimulated[0]
		}
	}
}
