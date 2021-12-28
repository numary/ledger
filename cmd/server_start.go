package cmd

import (
	"context"
	"github.com/numary/ledger/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"net"
	"net/http"
)

func NewServerStart() *cobra.Command {
	return &cobra.Command{
		Use: "start",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := NewContainer(
				fx.Invoke(func(h *api.API) error {
					listener, err := net.Listen("tcp", viper.GetString(serverHttpBindAddressFlag))
					if err != nil {
						return err
					}

					go http.Serve(listener, h)
					go func() {
						select {
						case <-cmd.Context().Done():
						}
						err := listener.Close()
						if err != nil {
							panic(err)
						}
					}()

					return nil
				}),
			)
			terminated := make(chan struct{})
			go func() {
				app.Run()
				close(terminated)
			}()
			select {
			case <-cmd.Context().Done():
				return app.Stop(context.Background())
			case <-terminated:
			}

			return nil
		},
	}
}