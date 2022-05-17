package pkg

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/coreos/go-semver/semver"
	"github.com/pterm/pterm"
	"github.com/x0f5c3/manic-go/pkg/downloader"
	"strings"
	"time"
)

type Releases = []Release

func GetReleases() (Releases, error) {
	resp, err := Get("https://api.github.com/repos/PowerShell/PowerShell/releases")
	if err != nil {
		return nil, err
	}
	var res Releases
	err = json.Unmarshal(resp, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type ParsedRelease struct {
	Version *semver.Version
	SHAFile *Asset
	Native  *Asset
}

type Downloaded struct {
	Version *semver.Version
	SHASum  string
	Data    []byte
}

func (p *ParsedRelease) Download() (*Downloaded, error) {
	resp, err := Get(p.SHAFile.BrowserDownloadUrl)
	if err != nil {
		return nil, err
	}
	var SHASum string
	scanner := bufio.NewScanner(strings.NewReader(string(resp)))
	for scanner.Scan() {
		t := scanner.Text()
		if strings.Contains(t, FileExt) {
			SHASum = strings.Split(t, " ")[0]
		}
	}
	l := p.Native.Size
	dl, err := downloader.New(p.Native.BrowserDownloadUrl, SHASum, HttpClient, &l)
	if err != nil {
		return nil, err
	}
	err = dl.Download(5, 5, true)
	if err != nil {
		return nil, err
	}
	return &Downloaded{
		Version: p.Version,
		SHASum:  SHASum,
		Data:    *dl.Data,
	}, nil
}

func AskForVersion(r Releases) (*Release, error) {
	sel := survey.Select{VimMode: true}
	var elems []string
	relsMap := make(map[string]*Release)
	for _, v := range r {
		relsMap[v.TagName] = &v
		elems = append(elems, v.TagName)
	}
	sel.Options = elems
	sel.Message = "Select a version"
	answer := 0
	err := survey.AskOne(&sel, &answer)
	if err != nil {
		return nil, err
	}
	v := r[answer]
	//if !ok {
	//	return nil, errors.New("no release found in map")
	//}
	return &v, nil
}

func (r *Release) Parse() (*ParsedRelease, error) {
	vers, err := semver.NewVersion(strings.Replace(r.TagName, "v", "", -1))
	if err != nil {
		return nil, err
	}
	res := ParsedRelease{Version: vers}
	for _, v := range r.Assets {
		if v.Name == "hashes.sha256" {
			res.SHAFile = &v
		} else if strings.Contains(v.Name, FileExt) {
			res.Native = &v
		}
	}
	if res.Native == nil {
		return nil, errors.New("no native asset found")
	}
	if res.SHAFile == nil {
		pterm.Error.Printf("Can't find hashes.sha256")
		pterm.Debug.Printf("Data: %v\n", r.Assets)
		return nil, errors.New("no hashes.sha256 found")
	}
	return &res, nil
}

type Asset struct {
	Url      string  `json:"url"`
	Id       int     `json:"id"`
	NodeId   string  `json:"node_id"`
	Name     string  `json:"name"`
	Label    *string `json:"label"`
	Uploader struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadUrl string    `json:"browser_download_url"`
}

type Release struct {
	Url       string `json:"url"`
	AssetsUrl string `json:"assets_url"`
	UploadUrl string `json:"upload_url"`
	HtmlUrl   string `json:"html_url"`
	Id        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeId          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []Asset   `json:"assets"`
	TarballUrl      string    `json:"tarball_url"`
	ZipballUrl      string    `json:"zipball_url"`
	Body            string    `json:"body"`
	DiscussionUrl   string    `json:"discussion_url,omitempty"`
	Reactions       struct {
		Url        string `json:"url"`
		TotalCount int    `json:"total_count"`
		Field3     int    `json:"+1"`
		Field4     int    `json:"-1"`
		Laugh      int    `json:"laugh"`
		Hooray     int    `json:"hooray"`
		Confused   int    `json:"confused"`
		Heart      int    `json:"heart"`
		Rocket     int    `json:"rocket"`
		Eyes       int    `json:"eyes"`
	} `json:"reactions,omitempty"`
	MentionsCount int `json:"mentions_count,omitempty"`
}
