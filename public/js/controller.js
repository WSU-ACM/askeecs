var askeecsControllers = angular.module('askeecsControllers', []);

askeecsControllers.controller('QuestionListCtrl', ['$scope', '$http',
	function ($scope, $http) {
		$http.get('data/questions.json').success(function(data) {
			$scope.questions = data;
		});
	}
]);

askeecsControllers.controller('QuestionAskCtrl', ['$scope', '$http', '$window', '$sce',
	function ($scope, $http, $window, $sce) {
		$scope.md2Html = function() {
			$scope.html = $window.marked($scope.markdown);
			$scope.htmlSafe = $sce.trustAsHtml($scope.html);
		}

		$scope.processForm = function () {
			console.log($scope.markdown);
			console.log($scope.tags);
			console.log($scope.title);

			$http({
				method: 'POST',
				url: '/question',
				data: $.param({title:$scope.title, body: $scope.markdown, tags: $scope.tags.split(' ')})
			}).success(function(data) {
				if (!data.success) {
					
				}
				else
				{

				}
			});
		}

	}
]);

askeecsControllers.controller('QuestionDetailCtrl', ['$scope', '$routeParams', '$http',
	function ($scope, $routeParams, $http) {
		$http.get('data/questions.json').success(function(data) {
			$scope.question = data[$routeParams.questionId];
		});
	}
]);
