package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Post struct {
	Title string
	Link  string
}

func main() {
	// Set environment variable to use a system font that supports Korean characters
	// Try a different font that might have better Korean character support
	os.Setenv("FYNE_FONT", "/System/Library/Fonts/Supplemental/AppleGothic.ttf")

	// Set additional environment variables for better Korean text support
	os.Setenv("FYNE_SCALE", "1.0")   // Ensure proper scaling
	os.Setenv("FYNE_THEME", "light") // Use light theme for better text visibility

	// Create a new Fyne application
	a := app.New()
	w := a.NewWindow("헬다이버즈 시리즈 갤러리 - 헬망호 탭")
	w.Resize(fyne.NewSize(800, 600))

	// Create a list to display posts
	postsList := widget.NewList(
		func() int { return 0 }, // Initial empty list
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel(""),
				widget.NewHyperlink("", nil),
				widget.NewSeparator(),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			// This will be updated when we have posts
		},
	)

	// Create a scroll container for the posts list
	scrollContainer := container.NewScroll(postsList)

	// Create a status label
	statusLabel := widget.NewLabel("게시글을 가져오려면 새로고침 버튼을 클릭하세요.")

	// Create a refresh button
	refreshButton := widget.NewButtonWithIcon("새로고침", theme.ViewRefreshIcon(), func() {
		fetchPosts(postsList, w, statusLabel)
	})

	// Create the main layout
	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("헬다이버즈 시리즈 갤러리 - 헬망호 탭 게시글 목록"),
			container.NewHBox(
				refreshButton,
				layout.NewSpacer(),
				statusLabel,
			),
		),
		nil, nil, nil,
		scrollContainer,
	)

	w.SetContent(content)

	// Initial fetch of posts
	go fetchPosts(postsList, w, statusLabel)

	// Show the window and run the app
	w.ShowAndRun()
}

func fetchPosts(postsList *widget.List, w fyne.Window, statusLabel *widget.Label) {
	// URL to scrape - DCInside helldiversseries gallery, 헬망호 tab
	urlStr := "https://gall.dcinside.com/mgallery/board/lists/?id=helldiversseries&sort_type=N&search_head=60&page=1"

	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Show loading dialog
	progress := dialog.NewProgress("로딩 중", "게시글을 가져오는 중입니다...", w)
	progress.Show()

	// Make HTTP request in a goroutine
	go func() {
		defer progress.Hide()

		// Make HTTP request
		resp, err := client.Get(urlStr)
		if err != nil {
			dialog.ShowError(fmt.Errorf("HTTP 요청 실패: %v", err), w)
			return
		}
		defer resp.Body.Close()

		// Check if the response status code is OK
		if resp.StatusCode != http.StatusOK {
			dialog.ShowError(fmt.Errorf("HTTP 요청 실패 (상태 코드: %d)", resp.StatusCode), w)
			return
		}

		// Parse HTML document
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			dialog.ShowError(fmt.Errorf("HTML 파싱 실패: %v", err), w)
			return
		}

		// Collect posts
		var posts []Post

		// Find and extract post titles and links
		doc.Find(".gall_list tbody tr.ub-content").Each(func(i int, s *goquery.Selection) {
			// Skip notices and advertisements
			isNotice := s.Find(".gall_num").Text() == "공지"
			isAd := s.HasClass("ub-content-up") || s.Find(".gall_subject").Text() == "AD"

			if !isNotice && !isAd {
				// Extract title
				titleElement := s.Find(".gall_tit > a:first-child")
				title := strings.TrimSpace(titleElement.Text())

				// Extract link
				link, exists := titleElement.Attr("href")
				if exists && !strings.HasPrefix(link, "javascript") {
					// If the link is relative, make it absolute
					if strings.HasPrefix(link, "/") {
						link = "https://gall.dcinside.com" + link
					}

					// Validate the link - should contain a post number
					if strings.Contains(link, "no=") && title != "" {
						posts = append(posts, Post{
							Title: title,
							Link:  link,
						})
					}
				}
			}
		})

		// Update the UI from the goroutine
		// In Fyne 2.4, we can directly update the UI from a goroutine
		postsList.Length = func() int {
			return len(posts)
		}

		postsList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
			container := obj.(*fyne.Container)
			titleLabel := container.Objects[0].(*widget.Label)
			linkHyperlink := container.Objects[1].(*widget.Hyperlink)

			post := posts[id]
			titleLabel.SetText(fmt.Sprintf("%d. %s", id+1, post.Title))

			// Check if the link is valid
			_, err := url.Parse(post.Link)
			if err != nil {
				linkHyperlink.SetText("링크 오류")
				linkHyperlink.SetURLFromString("")
				return
			}

			linkHyperlink.SetText("링크 열기")
			linkHyperlink.SetURLFromString(post.Link)
		}

		// Refresh the list
		postsList.Refresh()

		// Update status label
		statusLabel.SetText(fmt.Sprintf("총 %d개의 게시글을 찾았습니다.", len(posts)))
	}()
}
