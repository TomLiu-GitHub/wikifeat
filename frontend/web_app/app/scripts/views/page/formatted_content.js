/*
 * Licensed to Wikifeat under one or more contributor license agreements.
 * See the LICENSE.txt file distributed with this work for additional information
 * regarding copyright ownership.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright
 * notice, this list of conditions and the following disclaimer in the
 * documentation and/or other materials provided with the distribution.
 *  Neither the name of Wikifeat nor the names of its contributors may be used
 * to endorse or promote products derived from this software without
 * specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

define([
    'jquery',
    'underscore',
    'marionette',
    'backbone.radio'
], function($,_,Marionette,Radio){

    'use strict';

    return Marionette.ItemView.extend({
        template: _.template('<div id="formatted"></div>'),

        initialize: function(){
            this.pluginsStarted = Radio.channel('plugin').request('get:pluginsStarted');
            this.pluginContentViews = [];
        },

        onRender: function(){
            this.$("#formatted").html(this.model.get("content").formatted);

        },

        onShow: function(){
            var self = this;
            $.when(this.pluginsStarted).done(function(){
                    self.loadContentPlugins();
            });
        },

        loadContentPlugins: function(){
            var contentFields = this.$("#formatted").find("[data-plugin]");
            //var pvs = [];
            _.each(contentFields, function(field){
                var pluginName = $(field).data('plugin');
                var resourceId = $(field).data('id');
                console.log("PLUGIN: " + pluginName + ", ID: " + resourceId);
                var pg = window[pluginName];
                if(typeof pg !== 'undefined') {
                    try {
                        var contentView = pg.getContentView(field, resourceId);
                        this.pluginContentViews.push(contentView);
                        contentView.render();
                    }
                    catch (e) { //Bad Plugin! Bad!
                        // console.log(e);
                    }
                } else {
                    console.log("Plugin " + pluginName + " is undefined");
                }
            }.bind(this));
        },

        onDestroy: function(){
            _.each(this.pluginContentViews, function(cv){
                if(typeof cv.destroy !== 'undefined'){
                    //It's a marionette view, just call destroy
                    cv.destroy();
                } else {
                    //Assume a plain backbone view
                    cv.undelegateEvents();
                    cv.remove();
                    cv.unbind();
                }
            });
        }

    });

});