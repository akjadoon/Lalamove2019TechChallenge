package main

import (
	"context"
	"fmt"
	"strings"
	"log"
	"os"
	"bufio"
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	semver.Sort(releases)
	seen := make(map[int64]int64)

	for i := len(releases)-1; i >= 0; i-- {
		if releases[i].Compare(*minVersion) == 1 && releases[i].Major == minVersion.Major{
			prev, exists := seen[releases[i].Minor]
			if exists && prev >= releases[i].Patch{
				continue
			}
			versionSlice = append(versionSlice, releases[i])
			seen[releases[i].Minor]	= releases[i].Patch
		}
	}

	return versionSlice
}
//Extract Semvar info from a Gitgub Repository Release
func GetAllReleases(releases []*github.RepositoryRelease) []*semver.Version{
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
			versionString=strings.Split(versionString, "-")[0] //Remove alpha/beta postfix
		}
		allReleases[i] = semver.New(versionString)
	}
	return allReleases
}
// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	
	//Read File 
	path:= os.Args[1]
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err) 
	}
	scanner := bufio.NewScanner(f)
	scanner.Scan() // Skip Header Line

	//Parse File line
    for scanner.Scan() {
		line :=strings.Split(scanner.Text(),",")
		names :=strings.Split(line[0], "/")

		releases, _, err := client.Repositories.ListReleases(ctx, names[0], names[1], opt)
		if err != nil {
			log.Fatal(err) 
		}
		minVersion:= semver.New(line[1])

		allReleases := GetAllReleases(releases) 
		versionSlice := LatestVersions(allReleases, minVersion)
		fmt.Printf("latest versions of %s/%s: %s\n", names[0], names[1], versionSlice)
    }

}
