package main

import (
	"fmt"
	"os"
  "strings"
  "sync"
)

var startFolder = "./structure"
var itemToFind = "kissa"
var wg = new(sync.WaitGroup)

func getItemCount(filePath string) (int, error) {
  data, err := os.ReadFile(filePath)
  if err != nil {
    return 0, err
  }
  return strings.Count(string(data), itemToFind), nil
}

func countItems(itemChan chan string, countChan chan int) {
  for {
    filePath, more := <- itemChan
    if more {
      count, err := getItemCount(filePath)
      if err == nil {
        countChan <- count
      }
    } else {
      break
    }
  }
  close(countChan)
}

func getDirItems(dirName string, itemChan chan string) {
  dirItems, err := os.ReadDir(dirName)
  if err == nil {
    for _, dirItem := range dirItems {
      itemName := dirItem.Name()
      path := dirName + "/" + itemName
      if strings.Contains(itemName, ".") {
        itemChan <- path
      } else {
        wg.Add(1)
        go getDirItems(path, itemChan)
      }
    }
  }
  wg.Done()
}

func closeItemChan(itemChan chan string) {
  wg.Wait()
  close(itemChan)
}

func main() {
  countChan := make(chan int)
  itemChan := make(chan string)
  wg.Add(1)
  go closeItemChan(itemChan)
  go getDirItems(startFolder, itemChan)
  go countItems(itemChan, countChan)

  itemsFound := 0
  for {
    count, more := <- countChan
    if more {
      itemsFound += count
    } else {
      break
    }
  }
  fmt.Printf("Word '%s' found %d times.\n", itemToFind, itemsFound)
}
