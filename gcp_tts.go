package main

import (
	"context"
	b64 "encoding/base64"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func synthesizeSpeechRequest(message, langCode, voiceName, gender string) (string, error) {

	ctx := context.Background()
	ssmlVoiceGender := getVoiceGender(gender)

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	req := texttospeechpb.SynthesizeSpeechRequest{

		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: message},
		},

		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: langCode,
			Name:         voiceName,
			SsmlGender:   ssmlVoiceGender,
		},

		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return "", err
	}

	audioB64 := b64.StdEncoding.EncodeToString(resp.AudioContent)

	return audioB64, nil
}

func getVoiceGender(ssmlGender string) texttospeechpb.SsmlVoiceGender {

	if ssmlGender == "MALE" {
		return texttospeechpb.SsmlVoiceGender_MALE
	}

	if ssmlGender == "FEMALE" {
		return texttospeechpb.SsmlVoiceGender_FEMALE
	}

	if ssmlGender == "NEUTRAL" {
		return texttospeechpb.SsmlVoiceGender_NEUTRAL
	}

	return texttospeechpb.SsmlVoiceGender_SSML_VOICE_GENDER_UNSPECIFIED

}
