// @license
// Copyright (C) 2019  Valentino Medimorec
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"flag"
	"fmt"
	"github.com/valicm/rabbinator/cmd"
	"os"
)

var (
	consumer = flag.String("consumer", "", "Consumer tag, should be unique. Used for distinction between multiple consumers.")
	config   = flag.String("config", "", "Optional. Declare specific directory where config files are located. Etc. /var/www/my_directory")
)

func main() {

	flag.Parse()

	// Consumer flag is required.
	if *consumer == "" {
		flag.PrintDefaults()
		fmt.Println("Consumer flag is required. It is used to distinct multiple consumers for same queue, and utilizes yaml configuration with same naming")
		os.Exit(1)
	}

	// Initialize configuration setup and RabbitMQ connection.
	cmd.Initialize(*consumer, *config)

}
