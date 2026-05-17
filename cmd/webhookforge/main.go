package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Freedomwithin/webhookforge/internal/providers"
	"github.com/Freedomwithin/webhookforge/internal/proxy"
	"github.com/Freedomwithin/webhookforge/internal/replay"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "proxy":
		proxyCmd := flag.NewFlagSet("proxy", flag.ExitOnError)
		port := proxyCmd.Int("port", 9000, "Port to listen on")
		target := proxyCmd.String("target", "", "Target URL to forward to")
		secret := proxyCmd.String("secret", "", "Webhook signing secret")
		providerName := proxyCmd.String("provider", "stripe", "Webhook provider (stripe, paddle, shopify)")
		proxyCmd.Parse(os.Args[2:])

		if *target == "" || *secret == "" {
			fmt.Println("Error: --target and --secret are required for proxy mode")
			proxyCmd.Usage()
			os.Exit(1)
		}

		provider, err := getProvider(*providerName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		p := &proxy.ProxyServer{
			Port:     *port,
			Target:   *target,
			Secret:   *secret,
			Provider: provider,
		}
		if err := p.Start(); err != nil {
			fmt.Printf("Proxy error: %v\n", err)
			os.Exit(1)
		}

	case "replay":
		replayCmd := flag.NewFlagSet("replay", flag.ExitOnError)
		file := replayCmd.String("file", "", "JSON file containing the webhook payload")
		target := replayCmd.String("target", "", "Target URL to fire at")
		secret := replayCmd.String("secret", "", "Webhook signing secret")
		providerName := replayCmd.String("provider", "stripe", "Webhook provider (stripe, paddle, shopify)")
		replayCmd.Parse(os.Args[2:])

		if *file == "" || *target == "" || *secret == "" {
			fmt.Println("Error: --file, --target, and --secret are required for replay mode")
			replayCmd.Usage()
			os.Exit(1)
		}

		provider, err := getProvider(*providerName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		r := &replay.ReplayEngine{
			File:     *file,
			Target:   *target,
			Secret:   *secret,
			Provider: provider,
		}
		if err := r.Run(); err != nil {
			fmt.Printf("Replay error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func getProvider(name string) (providers.Provider, error) {
	switch name {
	case "stripe":
		return &providers.StripeProvider{}, nil
	case "paddle":
		return &providers.PaddleProvider{}, nil
	case "shopify":
		return &providers.ShopifyProvider{}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: webhookforge <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  proxy   - Start a re-signing proxy server")
	fmt.Println("  replay  - Re-sign and fire a saved JSON payload")
	fmt.Println("\nRun 'webhookforge <command> --help' for more information on a command.")
}
