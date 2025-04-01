# Go Concepts to Review

## Language Features
- [ ] Struct tags (e.g. `json:"id"`)
  - [Official Go Struct Tags](https://golang.org/ref/spec#Struct_types)
  - [Practical Go: Struct Tags](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go)
- [ ] Array vs Slice syntax
  - [Go Blog: Slices](https://blog.golang.org/slices-intro)
- [ ] Interface implementation
  - [Go Interface Types](https://golang.org/doc/effective_go#interfaces)
- [ ] Error handling patterns
  - [Error Handling in Go](https://go.dev/blog/error-handling-and-go)
- [ ] Pointer receivers vs value receivers
  - [Go Methods](https://golang.org/doc/effective_go#methods)

## Web Development
- [ ] Form handling in Go
  - [Go Web Examples: Forms](https://gowebexamples.com/forms/)
- [ ] Template syntax and patterns
  - [Go Templates Documentation](https://pkg.go.dev/text/template)
- [ ] HTTP request/response lifecycle
  - [Go HTTP Handler](https://golang.org/doc/articles/wiki/#tmp_3)
- [ ] Middleware patterns
  - [Go Web Examples: Middleware](https://gowebexamples.com/basic-middleware/)
- [ ] Javascript DOM manipulation in templates
  - [MDN: DOM Manipulation](https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Client-side_web_APIs/Manipulating_documents)
- [ ] Template security considerations
  - [Go Template Security](https://pkg.go.dev/html/template#hdr-Security)
- [ ] Form data parsing in Go
  - [Go Form Parsing](https://pkg.go.dev/net/http#Request.ParseForm)

## Project Patterns
- [ ] Project structure best practices
  - [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [ ] Dependency injection
  - [Go Dependency Injection](https://www.alexedwards.net/blog/organising-database-access)
- [ ] Error handling strategies
  - [Go Error Handling Best Practices](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- [ ] Testing patterns
  - [Go Testing Package](https://pkg.go.dev/testing)

## Debugging Skills
- [ ] Go logging best practices
  - [Go Logger Package](https://pkg.go.dev/log)
  - [Effective Logging in Go](https://www.digitalocean.com/community/tutorials/how-to-use-the-logger-package-in-go)
- [ ] Browser Developer Tools
  - [Chrome DevTools Network](https://developer.chrome.com/docs/devtools/network/)
- [ ] HTTP request/response debugging
  - [Go HTTP Debugging Guide](https://pkg.go.dev/net/http/httputil#DumpRequest)
- [ ] Template debugging
  - [Go Template Debugging](https://pkg.go.dev/text/template#hdr-Functions)

## Form Handling Lessons Learned
- [ ] PUT Request Form Data Handling
  - Key Issues:
    1. Browser handling of PUT requests differs from POST
    2. Form data encoding must be explicitly set
    3. Content-Type headers are crucial
  - Solution Components:
    - Form attributes: `enctype="application/x-www-form-urlencoded"`
    - Fetch headers: `'Content-Type': 'application/x-www-form-urlencoded'`
    - Data conversion: `new URLSearchParams(formData).toString()`
  - References:
    - [MDN: Using FormData](https://developer.mozilla.org/en-US/docs/Web/API/FormData/Using_FormData_Objects)
    - [Go Form Parsing](https://pkg.go.dev/net/http#Request.ParseForm)

## Current Questions
- How do browsers handle different HTTP methods with form data?
- What are the best practices for handling PUT/DELETE requests in Go web apps?
- When should we use different form encoding types?

## Resources
- [Go Template Documentation](https://pkg.go.dev/html/template)
- [MDN JavaScript in HTML](https://developer.mozilla.org/en-US/docs/Learn/HTML/Howto/Use_JavaScript_within_a_webpage)
- [MDN DOM Manipulation](https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Client-side_web_APIs/Manipulating_documents)
- [Go Form Processing](https://pkg.go.dev/net/http#Request.ParseForm)

## Database Concepts
- [ ] BoltDB Fundamentals
  - Key Concepts:
    - Key-value storage
    - B+ tree implementation
    - ACID transactions
    - Single-writer, multiple-reader design
  - References:
    - [BoltDB GitHub](https://github.com/boltdb/bolt)
    - [BoltDB Docs](https://pkg.go.dev/go.etcd.io/bbolt)

- [ ] Transaction Management
  - [ ] Read transactions
  - [ ] Write transactions
  - [ ] Transaction patterns in Go
  - References:
    - [BoltDB Transactions](https://pkg.go.dev/go.etcd.io/bbolt#DB.Begin)
    - [Go Database Patterns](https://www.alexedwards.net/blog/organising-database-access)

- [ ] Data Serialization
  - [ ] JSON encoding/decoding
  - [ ] Binary encoding options
  - [ ] Error handling in serialization
  - References:
    - [Go JSON Package](https://pkg.go.dev/encoding/json)
    - [Go Binary Encoding](https://pkg.go.dev/encoding/gob) 