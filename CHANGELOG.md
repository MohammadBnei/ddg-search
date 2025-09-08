# Changelog

## [0.6.2](https://github.com/MohammadBnei/ddg-search/compare/0.6.1...0.6.2) (2025-09-08)

## [0.6.1](https://github.com/MohammadBnei/ddg-search/compare/0.6.0...0.6.1) (2025-09-08)

# [0.6.0](https://github.com/MohammadBnei/ddg-search/compare/0.5.0...0.6.0) (2025-09-08)


### Features

* **main.go:** add support for setting Swagger host from the HOST_URL environment variable to customize swagger documentation host address. ([b8c540d](https://github.com/MohammadBnei/ddg-search/commit/b8c540d7929b2e6eb634a7b776ae338328e1d1ac))

# [0.5.0](https://github.com/MohammadBnei/ddg-search/compare/0.4.0...0.5.0) (2025-09-08)


### Features

* add more realistic user agent options ([c4cce1c](https://github.com/MohammadBnei/ddg-search/commit/c4cce1c15d598ac7e766ee8e15dbdbf0d9da701b))
* add randomized headers function and test suite ([865c98e](https://github.com/MohammadBnei/ddg-search/commit/865c98e192a5fb8281a3b62771ecb8f0236a9b1b))
* add Search method and enhance HTTP headers in SearchLimited ([b84e300](https://github.com/MohammadBnei/ddg-search/commit/b84e300f3b960e082d57306e9d781d03429a7949))

# [0.4.0](https://github.com/MohammadBnei/ddg-search/compare/0.3.0...0.4.0) (2025-07-04)


### Features

* add rate limiter to DuckDuckGoService ([1862811](https://github.com/MohammadBnei/ddg-search/commit/1862811e98cb3f4715af504fa9bf07a6b4e8550a))

# [0.3.0](https://github.com/MohammadBnei/ddg-search/compare/0.2.0...0.3.0) (2025-07-02)


### Bug Fixes

* debug mode is now just debug ([41cc09d](https://github.com/MohammadBnei/ddg-search/commit/41cc09d8e62763ce69aba2dce108bd994c6ef80b))


### Features

* **k8s:** updated service to allow connection from outside ([e388457](https://github.com/MohammadBnei/ddg-search/commit/e388457c168d75692e50631d711378dd1c59e7d5))

# [0.2.0](https://github.com/MohammadBnei/ddg-search/compare/0.1.0...0.2.0) (2025-04-08)


### Features

* Implement HTML scraping and Markdown conversion for search results ([b7190c1](https://github.com/MohammadBnei/ddg-search/commit/b7190c18ae026a449a196c62c3710ac404aaba38))

# [0.1.0](https://github.com/MohammadBnei/ddg-search/compare/0.0.4...0.1.0) (2025-04-07)


### Features

* Add dev command to Makefile using gowatch for development mode ([9441a56](https://github.com/MohammadBnei/ddg-search/commit/9441a56d197aa78e9ef02bf58011fe6e3a5351ed))
* Configure slog logger level based on debug mode ([abbf90c](https://github.com/MohammadBnei/ddg-search/commit/abbf90c2e2ddc41083012d1dbe52d115acc702a3))
* Enhance logging middleware with request/response details and ms duration ([ca2abe3](https://github.com/MohammadBnei/ddg-search/commit/ca2abe3cd431f7cca8820e8be3c5fb7dc31eaf2b))
* Implement logging middleware for debugging requests/responses ([1b16ab2](https://github.com/MohammadBnei/ddg-search/commit/1b16ab24342ddebc8b0c7bfaf0b9b1471a49d93f))
* Integrate Swagger for API documentation and add response types ([fa898ed](https://github.com/MohammadBnei/ddg-search/commit/fa898edf26b5ef9a57a4ada6a9c96821d075c1cc))

## [0.0.4](https://github.com/MohammadBnei/ddg-search/compare/0.0.3...0.0.4) (2025-03-10)


### Bug Fixes

* Resolve test failures by updating HTTP client and retry configuration handling ([237ce29](https://github.com/MohammadBnei/ddg-search/commit/237ce29ed2352e902a61cca39942b7f90b3a30eb))
* Resolve test failures in DuckDuckGo search client and service ([44daa7a](https://github.com/MohammadBnei/ddg-search/commit/44daa7a627242c35c3c93782115cfcfa846ae366))

## [0.0.3](https://github.com/MohammadBnei/ddg-search/compare/0.0.2...0.0.3) (2025-03-07)


### Bug Fixes

* **ci:** removed lint from tdocker build ci step ([21d8503](https://github.com/MohammadBnei/ddg-search/commit/21d85032c57c91c326aaf0b085a00a783773c07a))

## 0.0.2 (2025-03-07)


### Bug Fixes

* removed old handler ([1848ef7](https://github.com/MohammadBnei/ddg-search/commit/1848ef784c9017925a127b74e0f5f920e7d0eb63))
* removed secret ([d68eeba](https://github.com/MohammadBnei/ddg-search/commit/d68eeba9749bbf0045f44e0aeb9286e4d07a0143))
