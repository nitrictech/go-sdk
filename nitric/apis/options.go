package apis

type (
	ApiOption    = func(api *api)
	RouteOption  = func(route Route)
	MethodOption = func(mo *methodOptions)
)

type JwtSecurityRule struct {
	Issuer    string
	Audiences []string
}

type methodOptions struct {
	security         []OidcOptions
	securityDisabled bool
}

func WithMiddleware(middleware ApiMiddleware) ApiOption {
	return func(api *api) {
		api.middleware = middleware
	}
}

func OidcRule(name string, issuer string, audiences []string) SecurityOption {
	return func(scopes []string) OidcOptions {
		return OidcOptions{
			Name:      name,
			Issuer:    issuer,
			Audiences: audiences,
			Scopes:    scopes,
		}
	}
}

func WithSecurityJwtRule(name string, rule JwtSecurityRule) ApiOption {
	return func(api *api) {
		if api.securityRules == nil {
			api.securityRules = make(map[string]interface{})
		}

		api.securityRules[name] = rule
	}
}

func WithSecurity(oidcOptions OidcOptions) ApiOption {
	return func(api *api) {
		if api.security == nil {
			api.security = []OidcOptions{oidcOptions}
		} else {
			api.security = append(api.security, oidcOptions)
		}
	}
}

func WithPath(path string) ApiOption {
	return func(api *api) {
		api.path = path
	}
}

func WithNoMethodSecurity() MethodOption {
	return func(mo *methodOptions) {
		mo.securityDisabled = true
	}
}

func WithMethodSecurity(oidcOptions OidcOptions) MethodOption {
	return func(mo *methodOptions) {
		mo.securityDisabled = false
		if mo.security == nil {
			mo.security = []OidcOptions{oidcOptions}
		} else {
			mo.security = append(mo.security, oidcOptions)
		}
	}
}
