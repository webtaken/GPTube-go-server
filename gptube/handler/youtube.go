package handler

import (
	"fmt"
	"gptube/config"
	"gptube/database"
	"gptube/models"
	"gptube/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func YoutubePreAnalysisHandler(c *fiber.Ctx) error {
	var body models.YoutubePreAnalyzerReqBody

	if err := c.BodyParser(&body); err != nil {
		errorResp := models.YoutubePreAnalyzerRespBody{
			Err: fmt.Errorf("%v", err).Error(),
		}
		c.JSON(errorResp)
		return c.SendStatus(http.StatusInternalServerError)
	}

	if body.VideoID == "" {
		errorResp := models.YoutubePreAnalyzerRespBody{
			Err: fmt.Errorf("please provide a videoID").Error(),
		}
		c.JSON(errorResp)
		return c.SendStatus(http.StatusBadRequest)
	}

	videoData, err := services.CanProcessVideo(&body)
	if err != nil {
		errResp := models.YoutubePreAnalyzerRespBody{
			Err: fmt.Errorf("%v", err).Error(),
		}
		c.JSON(errResp)
		return c.SendStatus(http.StatusBadRequest)
	}

	maxNumCommentsRequireEmail, _ := strconv.Atoi(config.Config("YOUTUBE_MAX_COMMENTS_REQUIRE_EMAIL"))
	successResp := models.YoutubePreAnalyzerRespBody{
		VideoID:       body.VideoID,
		Snippet:       videoData.Items[0].Snippet,
		RequiresEmail: videoData.Items[0].Statistics.CommentCount > uint64(maxNumCommentsRequireEmail),
		NumOfComments: int(videoData.Items[0].Statistics.CommentCount),
	}
	c.JSON(successResp)
	return c.SendStatus(http.StatusOK)
}

func YoutubeAnalysisHandler(c *fiber.Ctx) error {
	var body models.YoutubeAnalyzerReqBody

	if err := c.BodyParser(&body); err != nil {
		errorResp := models.YoutubeAnalyzerRespBody{
			Err: fmt.Errorf("%v", err).Error(),
		}
		c.JSON(errorResp)
		return c.SendStatus(http.StatusInternalServerError)
	}

	// This means we haven´t received email hence is a short video so we do
	// all the logic here and send the response instantly to the client
	if body.Email == "" {
		results, err := services.Analyze(body)
		if err != nil {
			// Sending the error to the user
			errorResp := models.YoutubeAnalyzerRespBody{
				Err: fmt.Sprintf(
					"GPTube analysis for YT video %q failed 😔, try again later or contact us.",
					body.VideoTitle,
				),
			}
			c.JSON(errorResp)
			return c.SendStatus(http.StatusInternalServerError)
		}

		// sending the results to the user
		successResp := models.YoutubeAnalyzerRespBody{
			VideoID:    body.VideoID,
			VideoTitle: body.VideoTitle,
			Results:    results,
		}
		// Here we must save the results to FireStore //
		doc, err := database.AddYoutubeResult(&successResp)
		if err != nil {
			// Sending the e-mail error to the user
			log.Printf("error saving data to firebase: %v\n", err.Error())
		} else {
			// Saving the resultID into the result2Store var to send the email
			successResp.ResultsID = doc.ID
		}
		////////////////////////////////////////////////
		c.JSON(successResp)
		return c.SendStatus(http.StatusOK)
	}

	// This means we have received email hence this video is large so we do all
	// the logic in the server and send the result back to the email of the user
	// Adding lead email to temporal database
	go func() {
		results, err := services.Analyze(body)
		if err != nil {
			// Sending the e-mail error to the user
			subjectEmail := fmt.Sprintf(
				"GPTube analysis for YT video %q failed 😔",
				body.VideoTitle,
			)
			log.Printf("%v\n", err.Error())
			go services.SendYoutubeErrorTemplate(subjectEmail, []string{body.Email})
			return
		}
		// Here we must save the results to FireStore //
		results2Store := models.YoutubeAnalyzerRespBody{
			VideoID:    body.VideoID,
			VideoTitle: body.VideoTitle,
			Email:      body.Email,
			Err:        "",
			Results:    results,
		}
		doc, err := database.AddYoutubeResult(&results2Store)
		if err != nil {
			// Sending the e-mail error to the user
			subjectEmail := fmt.Sprintf(
				"GPTube analysis for YT video %q failed 😔",
				body.VideoTitle,
			)
			log.Printf("%v\n", err.Error())
			go services.SendYoutubeErrorTemplate(subjectEmail, []string{body.Email})
			return
		}
		// Saving the resultID into the result2Store var to send the email
		results2Store.ResultsID = doc.ID
		////////////////////////////////////////////////

		// Sending the e-mail to the user
		subjectEmail := fmt.Sprintf(
			"GPTube analysis for YT video %q ready 😺!",
			body.VideoTitle,
		)
		go services.SendYoutubeSuccessTemplate(
			results2Store, subjectEmail, []string{body.Email})
		fmt.Printf("Number of comments analyzed Bert: %d\n", results.BertResults.SuccessCount)
		fmt.Printf("Number of comments analyzed Roberta: %d\n", results.RobertaResults.SuccessCount)
	}()

	return c.SendStatus(http.StatusOK)
}
