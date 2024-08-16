package browser

import (
	"fmt"
	"os"
	"time"

	"github.com/playwright-community/playwright-go"
)

func GetBrowser(headless bool) (playwright.Browser, error) {
	// Use local chrome
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}
	err := playwright.Install(runOption)
	if err != nil {
		return nil, fmt.Errorf("could not install playwright dependencies: %v", err)
	}

	// Initialize Playwright
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start Playwright: %w", err)
	}
	// defer pw.Stop()
	// Launch a new browser instance in slow motion so that posthog can do its thing
	option := playwright.BrowserTypeLaunchOptions{
		Channel:  playwright.String("chrome"),
		Headless: playwright.Bool(headless),
		SlowMo:   playwright.Float(2000),
		Devtools: playwright.Bool(false),
	}
	return pw.Chromium.Launch(option)
}

func GetBrowserCustomResolver(headless bool, domain, serverIP string) (playwright.Browser, error) {
	// Use local chrome
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: true,
	}
	err := playwright.Install(runOption)
	if err != nil {
		return nil, fmt.Errorf("could not install playwright dependencies: %v", err)
	}

	// Initialize Playwright
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start Playwright: %w", err)
	}
	// defer pw.Stop()
	// Add a custom resolver for the domain. This is useful when the domain is not
	// reachable from the internet.
	// See https://peter.sh/experiments/chromium-command-line-switches/ for all flags
	resolverFlag := fmt.Sprintf("--host-resolver-rules=MAP %s %s", domain, serverIP)
	fmt.Printf("launching browser with a custom resolver: %s", resolverFlag)
	// Launch a new browser instance in slow motion so that posthog can do its thing
	option := playwright.BrowserTypeLaunchOptions{
		Channel:  playwright.String("chrome"),
		Headless: playwright.Bool(headless),
		SlowMo:   playwright.Float(2000),
		Devtools: playwright.Bool(false),
		Args: []string{
			resolverFlag,
			"--no-sandbox",
		},
	}
	return pw.Chromium.Launch(option)
}

func InjectPostHog(page playwright.Page, distinctID string) error {
	posthogSecret := os.Getenv("POSTHOG_SECRET")
	if posthogSecret == "" {
		return fmt.Errorf("POSTHOG_SECRET environment variable is not set")
	}

	// Inject PostHog snippet
	// For more PostHog settings see: https://posthog.com/docs/libraries/js
	posthogScript := fmt.Sprintf(`
	!function(t,e){var o,n,p,r;e.__SV||(window.posthog=e,e._i=[],e.init=function(i,s,a){function g(t,e){var o=e.split(".");2==o.length&&(t=t[o[0]],e=o[1]),t[e]=function(){t.push([e].concat(Array.prototype.slice.call(arguments,0)))}}(p=t.createElement("script")).type="text/javascript",p.async=!0,p.src=s.api_host.replace(".i.posthog.com","-assets.i.posthog.com")+"/static/array.js",(r=t.getElementsByTagName("script")[0]).parentNode.insertBefore(p,r);var u=e;for(void 0!==a?u=e[a]=[]:a="posthog",u.people=u.people||[],u.toString=function(t){var e="posthog";return"posthog"!==a&&(e+="."+a),t||(e+=" (stub)"),e},u.people.toString=function(){return u.toString(1)+".people (stub)"},o="capture identify alias people.set people.set_once set_config register register_once unregister opt_out_capturing has_opted_out_capturing opt_in_capturing reset isFeatureEnabled onFeatureFlags getFeatureFlag getFeatureFlagPayload reloadFeatureFlags group updateEarlyAccessFeatureEnrollment getEarlyAccessFeatures getActiveMatchingSurveys getSurveys getNextSurveyStep onSessionId setPersonProperties".split(" "),n=0;n<o.length;n++)g(u,o[n]);e._i.push([i,s,a])},e.__SV=1)}(document,window.posthog||[]);
posthog.init('%s',{api_host:'https://us.i.posthog.com'});
posthog.identify('%s');
	`, posthogSecret, distinctID)
	_, err := page.Evaluate(posthogScript)
	if err != nil {
		return fmt.Errorf("could not inject PostHog script: %v", err)
	}
	time.Sleep(2 * time.Second)

	return nil
}
