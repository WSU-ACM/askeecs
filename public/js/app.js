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
	}
]);

askeecsApp.run(function($rootScope, $location, AuthService, FlashService) {
	var routesThatRequireAuth = ['/ask', '/vote'];

	$rootScope.$on('$routeChangeStart', function (event, next, current) {
		FlashService.clear();
		if(_(routesThatRequireAuth).contains($location.path()) && !AuthService.isLoggedIn())
		{
			FlashService.show("Please login to continue");
			$location.path('/login');
		}
	});
});

askeecsApp.config(function($httpProvider) {
	var logsOutUserOn401 = function ($location, $q, SessionService, FlashService) {
		var success = function (res) {
			return res;
		}
		var error   = function (res) {
			if(res.status === 401) { // HTTP NotAuthorized
				SessionService.unset('authenticated')
				FlashService.show(res.data.flash);
				$location.path("/login");
				return $q.reject(res)
			} else {
				return $q.reject(res)
			}
		}

		return function(promise) {
			return promise.then(success, error)
		}
	}

	$httpProvider.responseInterceptors.push(logsOutUserOn401);
})

askeecsApp.factory("SessionService", function () {
	return {
		get: function (key) {
			return sessionStorage.getItem(key);
		},
		set: function (key, val) {
			return sessionStorage.setItem(key, val);
		},
		unset: function (key) {
			return sessionStorage.removeItem(key);
		}
	}
});

askeecsApp.factory("AuthService", ['$http', '$location', 'SessionService', 'FlashService',
	function($http, $location, SessionService, FlashService) {

		var cacheSession = function (user) {
			SessionService.set('authenticated', true);
			SessionService.set('user', user);
		}

		var uncacheSession = function () {
			SessionService.unset('authenticated');
			SessionService.unset('user');
		}

		var loginError = function (res) {
			FlashService.show(res.flash);
		}

		return {
			login: function (credentials) {
				var login = $http.post("/login", credentials);
				login.success(cacheSession);
				login.success(FlashService.clear);
				login.error(loginError);
				return login;
			},
			logout: function () {
				var logout =  $http.get("/logout");
				logout.success(uncacheSession);
				return logout;
			},
			isLoggedIn: function () {
				return SessionService.get('authenticated');
			},
			currentUser: function () {
				if ( this.isLoggedIn() )
				{
					return SessionService.get('user');
				}

				return null;
			}
		}
	}
]);

askeecsApp.factory("FlashService", function ($rootScope) {
	return {
		show: function (msg) {
			$rootScope.flash = msg
		},
		clear: function () {
			$rootScope.flash = ""
		}
	}
});

