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

askeecsApp.run(function($rootScope, $location, AuthService, FlashService, SessionService) {
	var routesThatRequireAuth = ['/ask'];

	$rootScope.authenticated = SessionService.get('authenticated');
	$rootScope.user = SessionService.get('user');

	$rootScope.$on('$routeChangeStart', function (event, next, current) {
		FlashService.clear()
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
				FlashService.show(res.data.Message);
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

askeecsApp.factory("AuthService", ['$rootScope', '$http', '$location', 'SessionService', 'FlashService',
	function($rootScope, $http, $location, SessionService, FlashService) {

		var cacheSession = function (user) {
			SessionService.set('authenticated', true);
			SessionService.set('user', user);
			$rootScope.authenticated = true;
			$rootScope.user = user;

		}

		var uncacheSession = function () {
			SessionService.unset('authenticated');
			SessionService.unset('user');
			$rootScope.authenticated = false;
			$rootScope.user = {};
		}

		var loginError = function (res) {
			FlashService.show(res.Message);
		}

		var protect = function (secret, salt) {
			var SHA256 = new Hashes.SHA256;
			return SHA256.hex(salt + SHA256.hex(secret));
		}

		var hash = function () {
			var s = ""
			var SHA256 = new Hashes.SHA256;

			for ( var i = 0; i < arguments.length; i++) {
				s += arguments[i];
			}

			return SHA256.hex(s);

		}

		return {
			login: function (credentials, fn) {

				// Friendly vars
				var u = credentials.Username;
				var p = credentials.Password;

				credentials.Username = "";
				credentials.Password = "";

				// Get a salt for this session
				$http.post("/register/salt", {"Username" : u})
					.success(function(user_salt) {
						$http.post("/salt", {"Username" : u})
							.success(function(session_salt) {

								// Produce the "Password" to send
								p = protect (u + p, user_salt.Salt);
								p = hash( p , session_salt)

								// Try to login
								var login = $http.post("/login", {"Username": u, "Password": p, "Salt": session_salt});

								login.success(cacheSession);
								login.success(FlashService.clear);
								login.error(loginError);

								if ( typeof fn === "function" )
									login.success(fn);
							}
						)
					}
				)
			},
			logout: function (fn) {
				var logout =  $http.post("/logout");
				logout.success(uncacheSession);

				if ( typeof fn === "function" )
					logout.success(fn);

			},
			register: function (credentials, fn) {

				// Friendly vars
				var u = credentials.Username;
				var p = credentials.Password;

				credentials.Username = "";
				credentials.Password = "";

				var s = protect(Date.now(), Math.random());

				// Produce the "Password" to send
				p = protect(u + p, s);

				var register = $http.post("/register", {"Username" : u, "Password" : p, "Salt" : s });

				if ( typeof fn === "function")
					register.success(fn);

			},
			isLoggedIn: function () {
				return SessionService.get('authenticated');
			},
			currentUser: function () {
				if ( this.isLoggedIn() )
				{
					return SessionService.get('user');
				}

				return {};
			}
		}
	}
]);

askeecsApp.factory("FlashService", function ($rootScope) {
	return {
		show: function (msg) {
			$rootScope.flashn = 1;
			$rootScope.flash = msg
		},
		clear: function () {
			if ( $rootScope.flashn-- == 0 )
				$rootScope.flash = ""
		}
	}
});

askeecsApp.factory('Questions', ['$http',
	function ($http) {

		var urlBase = '/q'
		var store	= []
		var f		= {};

		var p = function (data) {
			this.success = function (fn) {
				fn(data)
			}
		}

		f.List = function () {
			return $http.get(urlBase)
				.success(function (data) {
					store = data;
				});
		}

		f.Get = function (id, force) {
			if ( !force ) {
				for ( var i = 0; i < store.length; i++ )
				{
					if ( store[i].ID == id )
						return new p(store[i]); 
				}
			}
			
			return $http.get(urlBase + '/' + id);
		}

		f.Insert = function (item) {
			return $http.post(urlBase, item)
				.success(function(data) {
					store.push(data);
				});
		}

		f.Update = function (item) {
			return $http.put(urlBase + '/' + item.ID, item)
				.success(function (data) {
					for ( var i = 0; i < store.length; i++ )
					{
						if ( store[i].ID == id )
							return store[i] = data;
					}
				});
		}

		f.Delete = function (id) {
			return $http.delete(urlBase + '/' + id)
				.success(function (data) {
					for ( var i = 0; i < store.length; i++ )
					{
						if ( store[i].ID == id )
							return store.splice(i, 1);
					}
				})
		}
		
		return f;
	}
]);

askeecsApp.directive('askeecsLogout', function (AuthService) {
	return {
		restrict: 'A',
		 link: function(scope, element, attrs) {
			var evHandler = function(e) {
				e.preventDefault;
				AuthService.logout();
				return false;
			}

			element.on ? element.on('click', evHandler) : element.bind('click', evHandler);
		 }
	}
});

askeecsApp.directive('question', ['Questions',
	function (Questions) {
		function link ( scope, element, attributes ) {
			console.log("Generating question...", attributes.question);
			Questions.Get(attributes.question).success(function(data) {
				console.log(data)
			})
		}

		return {
			restruct: 'A',
			link: link
		}
	}
]);

askeecsApp.directive('comment', ['Questions',
	function (Questions) {
		function link ( scope, element, attributes ) {
			console.log("Generating comment...", attributes.question);
			Questions.Get(attributes.question).success(function(data) {
				console.log(data)
			})
		}

		return {
			restruct: 'A',
			link: link
		}
	}
]);

askeecsApp.filter('commentremark', function () {
	return function(input) {
		if(input === 0)
			return "at least enter 15 characters";
		else if(input < 15)
			return "" + 15 - input + " more to go..."
		else
			return 600 - input + " characters left"
		
	}
});
