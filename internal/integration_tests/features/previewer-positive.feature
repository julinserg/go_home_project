# file: features/previewer-positive.feature

# http://localhost:8088/
# http://previewer_service:8088/

Feature: Позитивные тесты превьюера изображений

	Scenario: Доступность сервиса превьюера
		When I send "GET" request to "http://previewer_service:8088/"
		Then The response code should be 200
		And The response should match text "This is my previewer!"
	
	Scenario: Получение превью изображения gopher_50x50.jpg 
		When I send "GET" request to "http://previewer_service:8088/fill/50/50/web/images/_gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_50x50.jpg"

	Scenario: Получение превью изображения gopher_500x500.jpg 
		When I send "GET" request to "http://previewer_service:8088/fill/500/500/web/images/_gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_500x500.jpg"

	Scenario: Получение превью изображения gopher_200x700.jpg 
		When I send "GET" request to "http://previewer_service:8088/fill/200/700/web/images/_gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_200x700.jpg"

	Scenario: Получение превью изображения gopher_256x126.jpg 
		When I send "GET" request to "http://previewer_service:8088/fill/256/126/web/images/_gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_256x126.jpg"

	Scenario: Получение превью изображения gopher_333x666.jpg 
		When I send "GET" request to "http://previewer_service:8088/fill/333/666/web/images/_gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_333x666.jpg"

	Scenario: Получение превью изображения gopher_1024x252.jpg 
		When I send "GET" request to "http://previewer_service:8088/fill/1024/252/web/images/_gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_1024x252.jpg"
	
	Scenario: Получение превью изображения gopher_2000x1000.jpg 
		When I send "GET" request to "http://previewer_service:8088/fill/2000/1000/web/images/_gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_2000x1000.jpg"