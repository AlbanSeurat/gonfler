package lzma


import (
	. "github.com/onsi/gomega"
	"testing"
	"os"
	"io/ioutil"
	"fmt"
)

func TestNewDecoder(t *testing.T) {
	RegisterTestingT(t)

	if file, err := os.Open("testdata/data.eos.l3.lzma") ; err == nil {
		defer file.Close()
		var header [13]byte
		if _,  err = file.Read(header[:]) ; err == nil {
			if props , err := NewProps(header[:]) ; err == nil {
				Expect(*props).Should(Equal(Props{ lc : 3, lp : 0, pb : 2, dicSize: 1048576 }))

				if fileContent, err := ioutil.ReadAll(file) ; err == nil {
					fmt.Println(fileContent)

					//TODO : manageDecodeProperly

				} else {
					t.Fatal(err)
				}

			} else {
				t.Fatal(err)
			}
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
}


