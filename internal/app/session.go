package app

import (
	"os"

	"github.com/Ruohao1/penta/internal/session"
	"github.com/Ruohao1/penta/internal/utils"
	"github.com/spf13/cobra"
)

func NewSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage sessions",
	}
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newUseCmd())
	// cmd.AddCommand(NewListCmd())
	return cmd
}

func newCreateCmd() *cobra.Command {
	var (
		name      string
		workspace string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new session",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := utils.LoggerFrom(cmd.Context())

			fileStore, err := session.NewFileStore()
			if err != nil {
				return err
			}

			if workspace == "" {
				workspace, err = os.Getwd()
				if err != nil {
					logger.Error().Err(err).Msg("Failed to get current directory")
					return err
				}
			}

			sessionID, err := utils.NewID()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to generate session ID")
				return err
			}

			session := session.Session{
				ID:        sessionID,
				Name:      name,
				Workspace: workspace,

				CreatedAt: utils.Now(),
				UpdatedAt: utils.Now(),
			}

			logger.Debug().Msgf("Creating session...")
			if err := fileStore.CreateSession(cmd.Context(), session); err != nil {
				if err == utils.ErrSessionExists {
					return nil
				}
				logger.Error().Err(err).Msg("Failed to create session")
				return err
			}
			logger.Info().Msgf("Session %s created successfully", name)

			if err := fileStore.SetCurrentSession(cmd.Context(), sessionID); err != nil {
				logger.Error().Err(err).Msg("Failed to set current session")
				return err
			}
			logger.Info().Msgf("Session %s set as current", name)

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Session name")
	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "Session workspace")

	cmd.MarkFlagRequired("name")

	return cmd
}

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := utils.LoggerFrom(cmd.Context())

			fileStore, err := session.NewFileStore()
			if err != nil {
				return err
			}

			// TODO: check the input
			logger.Debug().Msgf("Deleting session...")
			if err := fileStore.DeleteSession(cmd.Context(), args[0]); err != nil {
				logger.Error().Err(err).Msg("Failed to delete session")
				return err
			}
			logger.Info().Msgf("Session deleted successfully")

			lastUpdated, err := fileStore.LastUpdatedSession(cmd.Context())
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get last updated session")
			}
			if err := fileStore.SetCurrentSession(cmd.Context(), lastUpdated); err != nil {
				logger.Error().Err(err).Msg("Failed to set current session")
			}
			logger.Info().Msgf("Using last updated session %s...", lastUpdated)

			return nil
		},
	}
	return cmd
}

func newUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <id>",
		Short: "Use a session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := utils.LoggerFrom(cmd.Context())

			fileStore, err := session.NewFileStore()
			if err != nil {
				return err
			}
			// TODO: check the input
			if err := fileStore.SetCurrentSession(cmd.Context(), args[0]); err != nil {
				logger.Error().Err(err).Msgf("Failed to use session %s", args[0])
				return err
			}
			logger.Info().Msgf("Using session %s...", args[0])

			return nil
		},
	}
	return cmd
}
