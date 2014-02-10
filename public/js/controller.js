var askeecsControllers = angular.module('askeecsControllers', ['ngCookies']);

askeecsControllers.controller('QuestionListCtrl', ['$scope', '$http',
	function ($scope, $http) {
		$http.get('/q').success(function(data) {
			$scope.questions = data;
		});
	}
]);

askeecsControllers.controller('RegisterCtrl', ['$scope', '$http', '$location',
	function ($scope, $http, $location) {
		$scope.data = {}
		$scope.processForm = function () {
			console.log("GO!");
			if($scope.data.Password != $scope.data.cpassword) {
				console.log("Missed matched password");
				return;
			}

			delete $scope.data.cpassword;
			console.log($scope.data);
			$http({
				method: 'POST',
				url: '/register',
				data: $scope.data
			}).success(function(data) {
				$location.path("/login");	
			});
			
		}
	}
]);

askeecsControllers.controller('LoginCtrl', ['$rootScope', '$scope', '$http', '$cookies', '$location', 'AuthService',
	function ($rootScope, $scope, $http, $cookies, $location, AuthService) {
		$scope.credentials = { "Username": "", "Password": "" }
		$scope.processForm = function () {
			AuthService.login($scope.credentials).success(function () {
				$location.path('/questions');
			});
		}
	}
]);

askeecsControllers.controller('QuestionAskCtrl', ['$scope', '$http', '$window', '$sce', '$location',
	function ($scope, $http, $window, $sce, $location) {
		$scope.markdown="";
		$scope.title="";
		$scope.tags="";
		$scope.md2Html = function() {
			$scope.html = $window.marked($scope.markdown);
			$scope.htmlSafe = $sce.trustAsHtml($scope.html);
		}

		$scope.processForm = function () {
			console.log($scope.markdown);
			console.log($scope.tags);
			console.log($scope.title);
			delete $scope.errorMarkdown;
			delete $scope.errorTitle;
			delete $scope.errorTags;

			var err = false;

			if ($scope.markdown.length < 1)
			{
				$scope.errorMarkdown = "Your question must be 120 characters or more."
				err = true;
			}

			if ($scope.title.length == 0)
			{
				$scope.errorTitle = "You must enter a title."
				err = true;
			}

			if ($scope.tags.length == 0)
			{
				$scope.errorTags = "You must have at least one tag."
				err = true;
			}




			if (err) {
				return;
			}

			$http({
				method: 'POST',
				url: '/q',
				data: {Title:$scope.title, Body: $scope.markdown, Tags: $scope.tags.split(' ')}
			}).success(function(data) {
				console.log(data);
				$location.path("/questions/"+data);	
			});
		}

	}
]);

askeecsControllers.controller('QuestionDetailCtrl', ['$scope', '$routeParams', '$http',
	function ($scope, $routeParams, $http) {
		$http.get('/q/' + $routeParams.questionId).success(function(data) {
			$scope.question = data;
		});
	}
]);
