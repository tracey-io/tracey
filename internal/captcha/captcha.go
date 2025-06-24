package captcha

type Service struct {
	questionManager *QuestionManager
	tokenManager    *TokenManager
	powManager      *POWManager
}

func NewService(questionManager *QuestionManager, tokenManager *TokenManager, powManager *POWManager) *Service {
	return &Service{
		questionManager: questionManager,
		tokenManager:    tokenManager,
		powManager:      powManager,
	}
}

type Challenge struct {
	ID         string   `json:"id"`
	Prompt     string   `json:"prompt"`
	Options    []string `json:"options"`
	Token      string   `json:"token"`
	POWNonce   string   `json:"powNonce"`
	Difficulty int      `json:"difficulty"`
	Timestamp  int64    `json:"timestamp"`
}

func (s *Service) Generate() (*Challenge, error) {
	question := s.questionManager.Generate()

	token, err := s.tokenManager.SignAnswerToken(question)

	if err != nil {
		return nil, err
	}

	challenge := s.powManager.GenerateChallenge()

	return &Challenge{
		ID:         question.ID,
		Prompt:     question.Prompt,
		Options:    question.Options,
		Token:      token,
		POWNonce:   challenge.Nonce,
		Difficulty: challenge.Difficulty,
		Timestamp:  challenge.Timestamp,
	}, nil
}

func (s *Service) Verify(userAnswer, answerToken string, nonce string, counter int, timestamp int64) (string, error) {
	if err := s.powManager.Verify(nonce, counter, timestamp); err != nil {
		return "", err
	}

	isValid, err := s.tokenManager.VerifyAnswerToken(answerToken, userAnswer)

	if err != nil {
		return "", err
	}

	if !isValid {
		return "", ErrInvalidToken
	}

	passToken, err := s.tokenManager.SignPassToken()

	if err != nil {
		return "", err
	}

	return passToken, nil
}

func (s *Service) Validate(passToken string) (bool, error) {
	return s.tokenManager.VerifyPassToken(passToken)
}
