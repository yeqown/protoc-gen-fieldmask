package main

import (
	"fmt"
)

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║   FieldMask End-to-End Demonstration                          ║")
	fmt.Println("║   Client → Server → Client Flow                             ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")

	// Create server and client
	server := NewServer()
	client := NewClient(server)

	// ========================================================================
	// SCENARIO 1: FILTER Mode (Partial Response)
	// ========================================================================
	fmt.Println("\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  SCENARIO 1: FILTER MODE - Partial Response")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  In FILTER mode, only MARKED fields are returned.")
	fmt.Println("  Unmarked fields are cleared (set to zero/empty).")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	client.Scenario1_MobileClient()
	client.Scenario1_WebClient()
	client.Scenario1_NoFieldMask()

	// ========================================================================
	// SCENARIO 2: PRUNE Mode (Partial Update)
	// ========================================================================
	fmt.Println("\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  SCENARIO 2: PRUNE MODE - Partial Update")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  In PRUNE mode, MARKED fields are INCLUDED in update.")
	fmt.Println("  Only marked fields are updated.")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	client.Scenario2_UpdateProfileOnly()
	client.Scenario2_UpdateContactOnly()

	// ========================================================================
	// Summary
	// ========================================================================
	fmt.Println("\n")
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║   Summary                                                     ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════════╣")
	fmt.Println("║                                                               ║")
	fmt.Println("║   FILTER Mode (Partial Response):                            ║")
	fmt.Println("║   • Client marks fields to INCLUDE in response               ║")
	fmt.Println("║   • Server returns only marked fields                        ║")
	fmt.Println("║   • Use case: Reduce payload, hide sensitive data            ║")
	fmt.Println("║                                                               ║")
	fmt.Println("║   PRUNE Mode (Partial Update):                               ║")
	fmt.Println("║   • Client marks fields to INCLUDE in update                  ║")
	fmt.Println("║   • Server updates only marked fields                        ║")
	fmt.Println("║   • Use case: Incremental updates, prevent accidental changes ║")
	fmt.Println("║                                                               ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
}
