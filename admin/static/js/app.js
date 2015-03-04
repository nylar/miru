'use strict';

require.config({
    paths: {
        jquery: 'lib/jquery/jquery-min',
        underscore: 'lib/underscore/underscore-min',
        backbone: 'lib/backbone/backbone-min',
    }
});

requirejs(["app/main"]);