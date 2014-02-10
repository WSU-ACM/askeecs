var askeecsControllers = angular.module('askeecsControllers', []);

askeecsControllers.controller('QuestionListCtrl', ['$scope', '$http',
	function ($scope, $http) {
		$http.get('data/questions.json').success(function(data) {
			$scope.questions = data;
		});
	}
]);

askeecsControllers.controller('QuestionDetailCtrl', ['$scope', '$routeParams', '$http',
	function ($scope, $routeParams, $http) {
		$http.get('data/questions.json').success(function(data) {
			$scope.question = data[$routeParams.questionId];
		});
	}
]);
