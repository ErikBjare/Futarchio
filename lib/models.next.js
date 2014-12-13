this.SimpleSchema.debug = true;

var Schemas = {};
this.Schemas = Schemas;

Schemas.UserCreated = new SimpleSchema({
    createdBy: {
        type: String
    },
    createdAt: {
        type: Date
    }
});

Schemas.Post = new SimpleSchema([Schemas.UserCreated, {
    title: {
        type: String
    },
    description: {
        type: String
    }
}]);

Schemas.Poll = new SimpleSchema([Schemas.Post, {
    type: {
        type: String,
        allowedValues: ["YesNo", "YesNoRange"]
    }
}]);

Schemas.Statement = new SimpleSchema([Schemas.Post, {}]);

Schemas.Prediction = new SimpleSchema([Schemas.UserCreated, {
    credence: {
        type: Number
    },
    statement: {
        type: String
    }
}]);

// TODO: Separate Up/Down votes and votes on polls
Schemas.Vote = new SimpleSchema([Schemas.UserCreated, {
    type: {
        type: String,
        allowedValues: ["UpDown", "YesNo", "YesNoRange"]
    },
    value: {
        type: Number,
        decimal: true
    },
    post: {
        type: String
    }
}]);



class UserCreated {
    constructor() {
        this.createdBy = Meteor.userId();
        this.createdAt = new Date();
    }
}

class Post extends UserCreated {
    constructor(data) {
        super(data);
        this.title = data.title;
        this.description = data.description;
    }
}

this.Polls = new Mongo.Collection("polls");
this.Polls.attachSchema(Schemas.Poll);
class Poll extends Post {
    constructor(data) {
        super(data);
        this.type = data.type;
    }
}
this.Poll = Poll;

this.Statements = new Mongo.Collection("statements");
this.Statements.attachSchema(Schemas.Statement);
class Statement extends Post {
    constructor(data) {
        super(data);
    }
}
this.Statement = Statement;

this.Predictions = new Mongo.Collection("predictions");
this.Predictions.attachSchema(Schemas.Prediction);
class Prediction extends UserCreated {
    constructor(data) {
        super();
        this.credence = data.credence;
    }
}
this.Prediction = Prediction;

this.Votes = new Mongo.Collection("votes");
this.Votes.attachSchema(Schemas.Vote);
class Vote extends UserCreated {
    constructor(data) {
        super();
        this.type = data.type;
        this.value = Number(data.value);
        this.post = data.post;
    }
}
this.Vote = Vote;
