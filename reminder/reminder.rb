require 'sinatra'
require 'pry'
require 'json'

set :port, 8066
set :bind, '0.0.0.0'

set :mattermost_token, 'muj667ku6pyadye86zgxodsrfa'

get '/status' do
  "I'm fine!"
end

post '/' do
  return status 403 unless params['token'] == settings.mattermost_token
  content_type :json
  { text: params['token'] }.to_json
end
