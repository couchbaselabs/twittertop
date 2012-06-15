package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

var sampledata = `{"in_reply_to_status_id_str":null,"entities":{"user_mentions":[{"indices":[3,12],"id_str":"43630849","screen_name":"Annak808","name":"Anna Karyme Shuayre","id":43630849}],"urls":[],"hashtags":[]},"text":"RT @Annak808: I cry every time I watch RENT. Angel was meant to have a wonderful life, angels don't die.... :'(","created_at":"Thu Jun 14 20:02:03 +0000 2012","place":null,"truncated":false,"in_reply_to_user_id_str":null,"in_reply_to_screen_name":null,"id_str":"213360675165184000","in_reply_to_user_id":null,"retweeted":false,"source":"\u003Ca href=\"http:\/\/twitter.com\/#!\/download\/ipad\" rel=\"nofollow\"\u003ETwitter for iPad\u003C\/a\u003E","contributors":null,"retweet_count":0,"retweeted_status":{"in_reply_to_status_id_str":null,"entities":{"user_mentions":[],"urls":[],"hashtags":[]},"text":"I cry every time I watch RENT. Angel was meant to have a wonderful life, angels don't die.... :'(","created_at":"Thu Jun 14 19:54:46 +0000 2012","place":null,"truncated":false,"in_reply_to_user_id_str":null,"in_reply_to_screen_name":null,"id_str":"213358842401140736","in_reply_to_user_id":null,"retweeted":false,"source":"\u003Ca href=\"http:\/\/twitter.com\/#!\/download\/iphone\" rel=\"nofollow\"\u003ETwitter for iPhone\u003C\/a\u003E","contributors":null,"retweet_count":0,"coordinates":null,"geo":null,"user":{"show_all_inline_media":false,"lang":"en","created_at":"Sun May 31 02:49:16 +0000 2009","profile_sidebar_border_color":"5ED4DC","profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2088279639\/image_normal.jpg","default_profile_image":false,"time_zone":null,"url":"http:\/\/fullmoon808.tumblr.com","notifications":null,"profile_use_background_image":true,"favourites_count":207,"id_str":"43630849","following":null,"profile_text_color":"3C3940","description":"Little Monster. Alien. Music . Starbucks. Books. Klaine. Broadway. Loki's Army. If loving fashion is a crime, we plead guilty .","verified":false,"profile_background_image_url":"http:\/\/a0.twimg.com\/profile_background_images\/413025412\/Cherry.gif","location":"","profile_link_color":"0099B9","listed_count":1,"statuses_count":2198,"followers_count":123,"profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2088279639\/image_normal.jpg","screen_name":"Annak808","profile_background_color":"c0d8d8","protected":false,"default_profile":false,"contributors_enabled":false,"geo_enabled":false,"profile_background_tile":true,"profile_background_image_url_https":"https:\/\/si0.twimg.com\/profile_background_images\/413025412\/Cherry.gif","name":"Anna Karyme Shuayre","is_translator":false,"follow_request_sent":null,"profile_sidebar_fill_color":"95E8EC","id":43630849,"utc_offset":null,"friends_count":160},"favorited":false,"id":213358842401140736,"in_reply_to_status_id":null},"coordinates":null,"geo":null,"user":{"show_all_inline_media":false,"lang":"en","created_at":"Wed Oct 14 02:37:44 +0000 2009","profile_sidebar_border_color":"eb5ea2","profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2265141336\/image_normal.jpg","default_profile_image":false,"time_zone":"Mexico City","url":null,"notifications":null,"profile_use_background_image":true,"favourites_count":292,"id_str":"82256923","following":null,"profile_text_color":"3f0f99","description":"Cunning stunt. Stunning cunt.","verified":false,"profile_background_image_url":"http:\/\/a0.twimg.com\/profile_background_images\/46904348\/fondomiyavi.jpg","location":"Right behind you, don't turn!","profile_link_color":"b30083","listed_count":0,"statuses_count":21671,"followers_count":134,"profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2265141336\/image_normal.jpg","screen_name":"KeairaRiona","profile_background_color":"ebc1dc","protected":false,"default_profile":false,"contributors_enabled":false,"geo_enabled":false,"profile_background_tile":false,"profile_background_image_url_https":"https:\/\/si0.twimg.com\/profile_background_images\/46904348\/fondomiyavi.jpg","name":"Danny Killian","is_translator":false,"follow_request_sent":null,"profile_sidebar_fill_color":"c185c4","id":82256923,"utc_offset":-21600,"friends_count":209},"favorited":false,"id":213360675165184000,"in_reply_to_status_id":null}`

var deletes = `{"delete":{"status":{"user_id_str":"212934523","id_str":"99282795398053889","id":99282795398053889,"user_id":212934523}}}`

func TestParsing(t *testing.T) {
	r := strings.NewReader(sampledata)
	d := json.NewDecoder(r)
	tweet, err := parseNext(d)
	if err != nil {
		t.Fatalf("Error parsing value from stream: %v", err)
	}

	expected := Tweet{
		Text: `RT @Annak808: I cry every time I watch RENT. ` +
			`Angel was meant to have a wonderful life, angels don't die.... :'(`}
	expected.Sender.User = "KeairaRiona"
	expected.Sender.Name = "Danny Killian"

	if !reflect.DeepEqual(tweet, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, tweet)
	}
}
