package cmd

import (
	"context"
	"fmt"

	awsservice "github.com/HenryOwenz/ciselect/internal/aws"
	"github.com/spf13/cobra"
)

var (
	summary string
	profile string
	region  string
)

func init() {
	// Add profile and region flags to all commands that interact with AWS
	for _, cmd := range []*cobra.Command{approveCmd, rejectCmd, listCmd} {
		cmd.Flags().StringVarP(&profile, "profile", "p", "", "AWS profile to use (required)")
		cmd.Flags().StringVarP(&region, "region", "r", "", "AWS region to use (required)")
		cmd.MarkFlagRequired("profile")
		cmd.MarkFlagRequired("region")
	}

	approveCmd.Flags().StringVarP(&summary, "summary", "s", "", "Approval summary message")
	rootCmd.AddCommand(approveCmd, listCmd, rejectCmd)
}

var approveCmd = &cobra.Command{
	Use:   "approve [pipeline-name] [stage-name] [action-name]",
	Short: "Approve a manual approval action in CodePipeline",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		service, err := awsservice.NewService(ctx, profile, region)
		if err != nil {
			return fmt.Errorf("failed to create AWS service: %w", err)
		}

		// First, get the approval token
		approvals, err := service.ListPendingApprovals(ctx)
		if err != nil {
			return fmt.Errorf("failed to list approvals: %w", err)
		}

		var token string
		for _, approval := range approvals {
			if approval.PipelineName == args[0] &&
				approval.StageName == args[1] &&
				approval.ActionName == args[2] {
				token = approval.Token
				break
			}
		}

		if token == "" {
			return fmt.Errorf("no pending approval found for pipeline '%s' stage '%s' action '%s'", args[0], args[1], args[2])
		}

		return service.ApproveAction(ctx, args[0], args[1], args[2], token, summary)
	},
}

var rejectCmd = &cobra.Command{
	Use:   "reject [pipeline-name] [stage-name] [action-name]",
	Short: "Reject a manual approval action in CodePipeline",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		service, err := awsservice.NewService(ctx, profile, region)
		if err != nil {
			return fmt.Errorf("failed to create AWS service: %w", err)
		}

		// First, get the approval token
		approvals, err := service.ListPendingApprovals(ctx)
		if err != nil {
			return fmt.Errorf("failed to list approvals: %w", err)
		}

		var token string
		for _, approval := range approvals {
			if approval.PipelineName == args[0] &&
				approval.StageName == args[1] &&
				approval.ActionName == args[2] {
				token = approval.Token
				break
			}
		}

		if token == "" {
			return fmt.Errorf("no pending approval found for pipeline '%s' stage '%s' action '%s'", args[0], args[1], args[2])
		}

		return service.RejectAction(ctx, args[0], args[1], args[2], token, summary)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List pending manual approvals",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		service, err := awsservice.NewService(ctx, profile, region)
		if err != nil {
			return fmt.Errorf("failed to create AWS service: %w", err)
		}

		approvals, err := service.ListPendingApprovals(ctx)
		if err != nil {
			return fmt.Errorf("failed to list approvals: %w", err)
		}

		if len(approvals) == 0 {
			fmt.Println("No pending approvals found")
			return nil
		}

		fmt.Println("Pending Approvals:")
		for _, approval := range approvals {
			fmt.Printf("Pipeline: %s\n", approval.PipelineName)
			fmt.Printf("Stage: %s\n", approval.StageName)
			fmt.Printf("Action: %s\n", approval.ActionName)
			fmt.Printf("Token: %s\n", approval.Token)
			fmt.Println("---")
		}

		return nil
	},
}
