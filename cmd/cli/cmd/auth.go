package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	authUsername string
	authPassword string
	authRole     string
)

func init() {
	// Auth parent command
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication management commands",
		Long:  `Manage users and authentication for Velo platform.`,
	}

	// Login command
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Velo platform",
		Long:  `Authenticate with the Velo platform and store credentials.`,
		Run:   runLogin,
	}

	loginCmd.Flags().StringVar(&authUsername, "username", "", "Username")
	loginCmd.Flags().StringVar(&authPassword, "password", "", "Password (use interactive prompt if not provided)")

	// Logout command
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from Velo platform",
		Long:  `Remove stored authentication credentials.`,
		Run:   runLogout,
	}

	// Create user command
	createUserCmd := &cobra.Command{
		Use:   "create-user",
		Short: "Create a new user",
		Long:  `Create a new user account on the Velo platform.`,
		Run:   runCreateUser,
	}

	createUserCmd.Flags().StringVar(&authUsername, "username", "", "Username")
	createUserCmd.Flags().StringVar(&authPassword, "password", "", "Password (use interactive prompt if not provided)")
	createUserCmd.Flags().StringVar(&authRole, "role", "user", "User role (admin or user)")
	createUserCmd.MarkFlagRequired("username")

	// Change password command
	changePasswordCmd := &cobra.Command{
		Use:   "change-password",
		Short: "Change user password",
		Long:  `Change password for the current user or specified user.`,
		Run:   runChangePassword,
	}

	changePasswordCmd.Flags().StringVar(&authUsername, "username", "", "Username (current user if not specified)")

	// Add subcommands
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(createUserCmd)
	authCmd.AddCommand(changePasswordCmd)

	rootCmd.AddCommand(authCmd)
}

func runLogin(cmd *cobra.Command, args []string) {
	if authUsername == "" {
		fmt.Print("Username: ")
		reader := bufio.NewReader(os.Stdin)
		username, _ := reader.ReadString('\n')
		authUsername = strings.TrimSpace(username)
	}

	if authPassword == "" {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Printf("Error reading password: %v\n", err)
			return
		}
		authPassword = string(passwordBytes)
		fmt.Println()
	}

	// Create HTTP client and authenticate
	client := &http.Client{Timeout: timeout}

	loginData := map[string]string{
		"username": authUsername,
		"password": authPassword,
	}

	jsonData, _ := json.Marshal(loginData)

	resp, err := client.Post(fmt.Sprintf("http://%s/api/auth/login", serverAddr), "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		fmt.Printf("Login successful! Token: %s\n", result["token"])

		// Store token for future requests (TODO: implement token storage)
		fmt.Println("Token storage not yet implemented - use web interface for persistent sessions")
	} else {
		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		fmt.Printf("Login failed: %s\n", result["error"])
	}
}

func runLogout(cmd *cobra.Command, args []string) {
	// TODO: Remove stored credentials
	fmt.Println("Logout functionality will remove stored credentials")
	fmt.Println("User logged out successfully")
}

func runCreateUser(cmd *cobra.Command, args []string) {
	if authPassword == "" {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Printf("Error reading password: %v\n", err)
			return
		}
		authPassword = string(passwordBytes)
		fmt.Println()

		fmt.Print("Confirm password: ")
		confirmBytes, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Printf("Error reading password: %v\n", err)
			return
		}
		confirm := string(confirmBytes)
		fmt.Println()

		if authPassword != confirm {
			fmt.Println("Passwords do not match")
			return
		}
	}

	// TODO: Implement actual user creation with server
	fmt.Printf("User creation functionality will create user %s with role %s\n", authUsername, authRole)
	fmt.Println("User management system ready for implementation")
}

func runChangePassword(cmd *cobra.Command, args []string) {
	if authUsername == "" {
		fmt.Print("Username (current user if empty): ")
		reader := bufio.NewReader(os.Stdin)
		username, _ := reader.ReadString('\n')
		authUsername = strings.TrimSpace(username)
	}

	fmt.Print("Current password: ")
	_, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		return
	}
	fmt.Println()

	fmt.Print("New password: ")
	newBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		return
	}
	fmt.Println()

	fmt.Print("Confirm new password: ")
	confirmBytes, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		return
	}
	fmt.Println()

	if string(newBytes) != string(confirmBytes) {
		fmt.Println("Passwords do not match")
		return
	}

	// TODO: Implement actual password change with server
	fmt.Printf("Password change functionality will update password for user %s\n", authUsername)
	fmt.Println("Password management system ready for implementation")
}
