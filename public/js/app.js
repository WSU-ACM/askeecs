var askeecsApp = angular.module('askeecs', ['angularMoment', 'ngRoute', 'askeecsControllers'])

askeecsApp.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.
			when('/questions', {
				templateUrl: 'partials/question-list.html',
				controller: 'QuestionListCtrl'
			}).
			when('/questions/:questionId', {
				templateUrl: 'partials/question-detail.html',
				controller: 'QuestionDetailCtrl'
			}).
			otherwise({
				redirectTo: '/questions'
			});
	}]);
