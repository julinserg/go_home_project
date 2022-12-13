# file: features/previewer-cache.feature

# http://localhost:8088/
# http://previewer_service:8088/

Feature: Тесты превьюера изображений с проверкой работы кеша

	Scenario: Доступность сервиса превьюера
		When I send "GET" request to "http://previewer_service:8088/"
		Then The response code should be 200
		And The response should match text "This is my previewer!"

	Scenario: Очистка кеша сервиса превьюера
		When I send "GET" request to "http://previewer_service:8088/clearcache"
		Then The response code should be 200
		And The response should match text "Clear cache is done!"
	
	Scenario: Получение превью изображения gopher_50x50.jpg с удаленного сервера
		When I send "GET" request to "http://previewer_service:8088/fill/50/50/web/images/gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_50x50.jpg"
		And Image get from remote server

	Scenario: Получение превью изображения gopher_50x50.jpg из кеша
		When I send "GET" request to "http://previewer_service:8088/fill/50/50/web/images/gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_50x50.jpg"
		And Image get from cache
	
	Scenario: Получение превью изображения gopher_50x50.jpg из кеша
		When I send "GET" request to "http://previewer_service:8088/fill/50/50/web/images/gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_50x50.jpg"
		And Image get from cache

	Scenario: Получение превью изображения gopher_200x700.jpg с удаленного сервера
		When I send "GET" request to "http://previewer_service:8088/fill/200/700/web/images/gopher_original_1024x504.jpg"
		Then The response code should be 200
		And The response equivalent image "./etalon_images/gopher_200x700.jpg"