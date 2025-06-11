# 헬다이버즈 시리즈 갤러리 - 헬망호 탭 스크래퍼

디시인사이드 헬다이버즈 시리즈 갤러리의 헬망호 탭에서 게시글을 가져와 표시하는 Fyne GUI 애플리케이션입니다.

## 기능

- 헬다이버즈 시리즈 갤러리 헬망호 탭에서 게시글 목록 가져오기
- 한글 지원을 위한 폰트 설정
- 게시글 제목과 링크를 GUI로 표시
- 새로고침 버튼으로 최신 게시글 업데이트

## 사용 기술

- **Go**: 메인 프로그래밍 언어
- **Fyne v2**: GUI 프레임워크
- **goquery**: HTML 파싱 및 웹 스크래핑

## 실행 방법

```bash
go run main.go
```

## 빌드

```bash
go build -o helldiver-scraper main.go
```

## 의존성

- `fyne.io/fyne/v2`: GUI 프레임워크
- `github.com/PuerkitoBio/goquery`: HTML 파싱

## 주요 특징

- 한글 텍스트 지원을 위한 AppleGothic 폰트 사용
- 공지사항 및 광고 게시글 필터링
- 클릭 가능한 하이퍼링크로 원본 게시글 접근
- 로딩 프로그레스 바와 상태 표시