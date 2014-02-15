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
		var dataMaster = {"Username" : "", "Password" : ""}
		$scope.data = {}
		$scope.processForm = function () {
			if($scope.data.Password != $scope.data.cpassword) {
				console.log("Missed matched password");
				return;
			}

			delete $scope.data.cpassword;


			// Generate a SHA256 Hasher
			var SHA256 = new Hashes.SHA256;

			// Friendly vars
			var u = $scope.data.Username;
			var p = $scope.data.Password;
			var s = "" + Date.now() % Math.random();
				s = SHA256.hex(s);

			// Reset the scope
			$scope.data = dataMaster;

			p = SHA256.hex(s + SHA256.hex(u + ":" + p));

			$http({
				method: 'GET',
				url: '/register',
				data: {"Username" : u, "Password" : p, "Salt" : s }
			}).success(function(data) {
				$location.path("/login");
			});
			
		}
	}
]);

askeecsControllers.controller('LoginCtrl', ['$scope', '$http', '$cookies', '$location', 'AuthService',
	function ($scope, $http, $cookies, $location, AuthService) {
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
			var src = $scope.markdown || ""
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

			if ($scope.markdown.length < 50)
			{
				$scope.errorMarkdown = "Your question must be 50 characters or more."
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

askeecsControllers.controller('QuestionDetailCtrl', ['$scope', '$routeParams', '$http', '$window', '$sce',
	function ($scope, $routeParams, $http, $window, $sce) {
		$scope.comment = { "Body" : "" }
		$scope.response = { "Body" : "" }

		$http.get('/q/' + $routeParams.questionId).success(function(data) {
			$scope.question = data;
			console.log(data)
		});

		$scope.voteUp = function () {
			$http({
				method: 'GET',
				url: '/q/' + $scope.question.ID + '/vote/up',
				data: {}
			}).success(function(data) {
				$scope.question.Upvotes = data.Upvotes
			});
		}

		$scope.voteDown = function () {
			$http({
				method: 'GET',
				url: '/q/' + $scope.question.ID + '/vote/down',
				data: {}
			}).success(function(data) {
				$scope.question.Downvotes = data.Downvotes
			});
		}

		$scope.markdown="";
		$scope.md2Html = function() {
			var src = $scope.response.Body || ""
			$scope.html = $window.marked(src);
			$scope.htmlSafe = $sce.trustAsHtml($scope.html);
		}

		$scope.processComment = function () {
			delete $scope.errorComment;

			var err = false;

			if ( $scope.comment.Body.length < 15 )
			{
				$scope.errorComment = "Your comment must be at least 15 characters"
				err = true;
			}

			if (err) return;

			$http({
				method: 'post',
				url: '/q/' + $scope.question.ID + '/comment/',
				data: $scope.comment
			}).success(function(data) {
				delete $scope.scomment;
				$scope.question.Comments.push(data);
			});
		}

		$scope.processForm = function () {
			console.log($scope.response.Body);
			delete $scope.errorMarkdown;

			var err = false;

			if ($scope.response.Body.length < 50)
			{
				$scope.errorMarkdown = "Your response must be 50 characters or more."
				err = true;
			}


			if (err) {
				return;
			}

			$http({
				method: 'post',
				url: '/q/' + $scope.question.ID + '/response/',
				data: $scope.response
			}).success(function(data) {
				$scope.question.Responses.push(data);
			});
		}
	}
]);
