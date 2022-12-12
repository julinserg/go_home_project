# file: features/previewer-negative.feature

# http://localhost:8088/
# http://previewer_service:8088/

Feature: Негативные тесты превьюера изображений

	Scenario: Доступность сервиса превьюера
		When I send "GET" request to "http://localhost:8088/"
		Then The response code should be 200
		And The response should match text "This is my previewer!"
	
	Scenario: Получение превью изображения для несуществующего файла 
		When I send "GET" request to "http://localhost:8088/fill/50/50/localhost/images/1.jpg"
		Then The response code should be 404
	
	Scenario: Получение превью изображения для файла не являющегося изображением 
		When I send "GET" request to "http://localhost:8088/fill/50/50/localhost/index.html"
		Then The response code should be 500
		And The response should match text "invalid JPEG format: missing SOI marker"

	Scenario: Получение превью изображения от несуществующего сервера 
		When I send "GET" request to "http://localhost:8088/fill/50/50/myimageserver/images/_gopher_original_1024x504.jpg"
		Then The response code should be 400
		And The response should match text
		"""
Get "http://myimageserver/images/_gopher_original_1024x504.jpg": dial tcp: lookup myimageserver: Temporary failure in name resolution
		"""

	Scenario: Получение превью изображения c неправильно заданной шириной 
		When I send "GET" request to "http://localhost:8088/fill/abc/50/localhost/images/_gopher_original_1024x504.jpg"
		Then The response code should be 400
		And The response should match text "width in URL not integer"

	Scenario: Получение превью изображения c неправильно заданной высотой 
		When I send "GET" request to "http://localhost:8088/fill/50/abc/localhost/images/_gopher_original_1024x504.jpg"
		Then The response code should be 400
		And The response should match text "height in URL not integer"
	
	Scenario: Получение превью изображения c неправильно количеством аргументов (нет пути до изображения) 
		When I send "GET" request to "http://localhost:8088/fill/50/50"
		Then The response code should be 400
		And The response should match text "not set width, height or image path in URL"

	Scenario: Получение превью изображения c неправильно количеством аргументов (нет высоты или ширины) 
		When I send "GET" request to "http://localhost:8088/fill/50"
		Then The response code should be 400
		And The response should match text "not set width, height or image path in URL"

	Scenario: Получение превью изображения c неправильно количеством аргументов (нет параметров) 
		When I send "GET" request to "http://localhost:8088/fill"
		Then The response code should be 400
		And The response should match text "not set width, height or image path in URL"
	
	Scenario: Получение превью изображения c неправильно количеством аргументов (слишком много параметров) 
		When I send "GET" request to "http://localhost:8088/fill/50/50/50/localhost/images/_gopher_original_1024x504.jpg"
		Then The response code should be 400
		And The response should match text
		"""
Get "http://50/localhost/images/_gopher_original_1024x504.jpg": dial tcp 0.0.0.50:80: connect: network is unreachable
		"""
	
	Scenario: Получение превью изображения со слишком большой шириной (width > 3840 || height > 2160) 
		When I send "GET" request to "http://localhost:8088/fill/4000/2160/localhost/images/_gopher_original_1024x504.jpg"
		Then The response code should be 400
		And The response should match text "width or height is very large"