var askeecsApp = angular.module('askeecs', ['angularMoment', 'ngRoute', 'askeecsControllers', 'ngCookies'])

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
			when('/ask', {
				templateUrl: 'partials/question-post.html',
				controller: 'QuestionAskCtrl'
			}).
			when('/register', {
				templateUrl: 'partials/register.html',
				controller: 'RegisterCtrl'
			}).
			when('/login', {
				templateUrl: 'partials/login.html',
				controller: 'LoginCtrl'
			}).
			otherwise({
				redirectTo: '/questions'
			});
	}]);

