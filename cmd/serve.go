/*
Copyright © 2021 Patryk Kalinowski <patryk@kalinowski.dev>
This file is part of the Biscuit API.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/finebiscuit/api/config"
	"github.com/finebiscuit/api/graph"
	"github.com/finebiscuit/api/graph/generated"
	"github.com/finebiscuit/api/sqldb"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the Biscuit API",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New(viper.GetViper())
		resolver, err := graph.NewResolver(cfg, sqldb.NewBackend())
		if err != nil {
			log.Fatal(err)
		}

		router := chi.NewRouter()

		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

		router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
		router.Handle("/graphql", srv)

		log.Printf("connect to http://localhost:%d/ for GraphQL playground", cfg.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serveCmd.Flags().Uint16P("port", "p", 8080, "Port to listen on")
	serveCmd.Flags().String("database-type", "sqlite3", "Database engine")
	serveCmd.Flags().String("database-source", "biscuit.db", "Database source")

	viper.BindPFlags(serveCmd.Flags())
}
