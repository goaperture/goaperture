gen:
	go run main.go --app app --routes api/routes --output api


prod:
	go build -o dist/goaperture main.go && mv dist/goaperture ~/go/bin