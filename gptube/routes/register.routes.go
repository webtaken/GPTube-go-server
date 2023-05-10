package routes

// func RegisterHandler(c *fiber.Ctx) error {
// 	var registerReq models.RegisterReq

// 	if err := c.BodyParser(registerReq); err != nil {
// 		ErrorResponse := models.RegisterResp{
// 			Err: fmt.Errorf("%v", err).Error(),
// 		}
// 		c.JSON(ErrorResponse)
// 		return c.SendStatus(http.StatusBadRequest)
// 	}

// 	if registerReq.Email == "" {
// 		ErrorResponse := models.RegisterResp{
// 			Err: fmt.Errorf("please provide an email").Error(),
// 		}
// 		data, err := json.Marshal(ErrorResponse)
// 		w.WriteHeader(http.StatusBadRequest)
// 		if err != nil {
// 			log.Printf("JSON marshaling failed: %s", err)
// 			return
// 		}
// 		w.Write(data)
// 		return
// 	}

// 	// Adding the lead
// 	err := firebase_services.AddLead(registerReq.Email)
// 	if err != nil {
// 		ErrorResponse := models.RegisterResp{
// 			Err: fmt.Errorf("couldn't add your email").Error(),
// 		}
// 		data, err := json.Marshal(ErrorResponse)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		if err != nil {
// 			log.Printf("JSON marshaling failed: %s", err)
// 			return
// 		}
// 		w.Write(data)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// }
