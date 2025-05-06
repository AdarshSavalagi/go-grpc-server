// Package middleware provides various middleware utilities, including request profiling.
package middleware

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/gin-gonic/gin"
)

// ProfilingMiddleware captures CPU profiling data for each incoming request.
//
// It starts profiling at the beginning of the request and stops profiling once the request is completed.
// The profiling data is saved to a file named "cpu.prof" for further analysis.
//
// Returns:
// - gin.HandlerFunc: A Gin middleware function.
func ProfilingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a file to store CPU profiling data
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}

		// Start CPU profiling
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile() // Ensure profiling stops when the request completes

		// Continue processing the request
		c.Next()
	}
}
