SimpleSchema.debug = true;

var Schemas = {};

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



UserCreated = function() {
    obj = {};
    obj.createdBy = Meteor.userId();
    obj.createdAt = new Date();
    return obj;
};

Post = function(data) {
    obj = UserCreated();
    obj.title = data.title;
    obj.description = data.description;
    return obj;
};

Polls = new Mongo.Collection("polls");
Polls.attachSchema(Schemas.Poll);
Poll = function(data) {
    obj = Post(data);
    obj.type = data.type;
    return obj;
};

Statements = new Mongo.Collection("statements");
Statements.attachSchema(Schemas.Statement);
Statement = function(data) {
    obj = Post(data);
    return obj;
};

Predictions = new Mongo.Collection("predictions");
Predictions.attachSchema(Schemas.Prediction);
Prediction = function(data) {
    obj = UserCreated();
    obj.credence = data.credence;
    obj.statement = data.statement;
    return obj;
};

Votes = new Mongo.Collection("votes");
Votes.attachSchema(Schemas.Vote);
Vote = function(data) {
    obj = UserCreated();
    obj.type = data.type;
    obj.value = Number(data.value);
    obj.post = data.post;
    return obj;
};
