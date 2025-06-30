/*
Copyright © 2025 Alireza Arzehgar <alirezaarzehgar82@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	Version     string
	updateParam struct {
		BinDir string
	}
)

type repoLatestVersion struct {
	Name   string `json:"name"`
	Assets []struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	} `json:"assets"`
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update tool to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		// Disable cursor
		fmt.Print("\x1b[?25l")
		defer fmt.Print("\x1b[?25h")

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			fmt.Print("\x1b[?25h")
			os.Exit(0)
		}()

		res, err := http.Get("https://api.github.com/repos/alirezaarzehgar/iflashc/releases/latest")
		if err != nil {
			log.Fatalf("failed to get latest version from github: %s", err)
		}

		rlv := repoLatestVersion{}
		err = json.NewDecoder(res.Body).Decode(&rlv)
		if err != nil {
			log.Fatalf("failed to unmarshal github response: %s", err)
		}

		if Version == rlv.Name {
			fmt.Println("No need to update!", "version", Version)
			return
		}

		outputFile, err := os.OpenFile("/tmp/iflashc", os.O_CREATE|os.O_WRONLY, os.FileMode(0755))
		if err != nil {
			log.Fatalf("failed to create binary file: %s", err)
		}
		defer func() {
			outputFile.Close()
			os.Remove("/tmp/iflashc")
		}()

		fmt.Println("downloading", rlv.Name)
		res, err = http.Get("https://github.com/alirezaarzehgar/iflashc/releases/latest/download/iflashc")
		if err != nil {
			log.Fatalf("failed to download latest version: %s", err)
		}
		defer res.Body.Close()

		var totalSize int64
		for _, asset := range rlv.Assets {
			if asset.Name == "iflashc" {
				totalSize = asset.Size
				break
			}
		}

		done := make(chan any)
		go func() {
			if totalSize == 0 {
				log.Println("no progress bar. total reported update size is zero!")
				return
			}

			for {
				time.Sleep(time.Second / 3)
				offset, err := outputFile.Seek(0, io.SeekCurrent)
				if err != nil {
					continue
				}
				if offset > 0 {
					progress := float64(offset) / float64(totalSize) * 50
					mb := float64(offset) / 1024 / 1024
					visualPerc := strings.Repeat("█", int(progress)) + strings.Repeat(" ", 50-int(progress))

					fmt.Printf("\r %.2fMB - %.2f%% [%s]", mb, progress*2, visualPerc)
				}
				if offset == totalSize {
					fmt.Println()
					done <- struct{}{}
					break
				}
			}
		}()

		n, err := io.Copy(outputFile, res.Body)
		if err != nil {
			log.Fatalf("failed to write on destination file: %s", err)
		}

		<-done
		fmt.Printf("iflashc updated successfully: %.2fMB\n", float64(n)/1024/1024)

		c := exec.Command("mv", "/tmp/iflashc", updateParam.BinDir)
		if err := c.Run(); err != nil {
			log.Fatalf("failed to mv downloaded update to destination: %s", err)
		}
	},
}

func init() {
	updateCmd.PersistentFlags().StringVar(&updateParam.BinDir, "dir", "/bin", "output directory for downloading binary")
	rootCmd.AddCommand(updateCmd)
}
