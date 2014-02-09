##Idea for Question object

    {
		"ID":"somehexvalue",
		"Title":"A fun question",
		"Author":"Jeromy",
		"Tags":["C","Meta","Fail"],
		"Score":11,
		"Timestamp":"12/14/14",
		"Body":"I dont know how to program, pls help",
		"Responses": [
			{
				"ID":"tltd",
				"Author":"TravisLane",
				"Timestamp":"12/15/14",
				"Score":45,
				"Body":"Noob, go to class more."
			}
		]
	}

##API Layout

###/question
POST - Post a new question, responds with link to page (or just id of question)
GET '/:id' - Get a question by id
POST '/:id/respond' - Send response to question
POST '/:id/up' - Upvote question
