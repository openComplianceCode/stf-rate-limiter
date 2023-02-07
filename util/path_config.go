package util

import (
	"io/ioutil"
	"log"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

type PathArray []*struct {
	Path        string `yaml:"path"`
	PathType    string `yaml:"pathType"`
	Consumption int    `yaml:"consumption"`
}

type PathConfig struct {
	Paths PathArray
}

func (ps PathArray) Len() int {
	return len(ps)
}

func (ps PathArray) Less(i, j int) bool {
	if ps[i].PathType == "Exact" {
		return ps[j].PathType == "Prefix" || len(ps[i].Path) > len(ps[j].Path)
	} else {
		return (ps[j].PathType == "Prefix") && len(ps[i].Path) > len(ps[j].Path)
	}
}

func (ps PathArray) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func ReadPathConfig() (*PathConfig, error) {
	var pathConfig PathConfig
	cpath, _ := os.Getwd()
	log.Println(cpath)

	file, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		file, err = ioutil.ReadFile("./conf/config.yaml")
	}
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &pathConfig)
	if err != nil {
		log.Fatalln("There is config file, but the format is wrong", err)
		return nil, err
	}

	if pathConfig.Paths != nil {
		for _, pathRule := range pathConfig.Paths {
			if pathRule.PathType == "" {
				pathRule.PathType = "Prefix"
			}
			if pathRule.Consumption == 0 {
				pathRule.Consumption = 1
			}
		}
	}

	sort.Sort(pathConfig.Paths)
	return &pathConfig, nil
}
