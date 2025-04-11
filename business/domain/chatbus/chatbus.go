package chatbus

import (
	"context"

	"github.com/google/uuid"
)

type Business struct {
}

func (b *Business) Conversation(ctx context.Context, con Conversation) (Conversation, error) {
	c := Conversation{}

	c.ID = uuid.New()
	c.ParentMessageID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

	m := Message{}
	m.ID = uuid.New()
	m.Role = RoleAssistant
	m.Content = "The sky appears blue to our eyes because of a phenomenon called Rayleigh scattering, named after the British physicist Lord Rayleigh, who first described it in the late 19th century.\n\nHere's what happens:\n\n1. When sunlight enters Earth's atmosphere, it encounters tiny molecules of gases such as nitrogen (N2) and oxygen (O2). These molecules are much smaller than the wavelength of visible light.\n2. The shorter (blue) wavelengths of light are scattered more than the longer (red) wavelengths by these small molecules. This is because the smaller molecules are more effective at scattering the shorter wavelengths due to their size.\n3. As a result, when we look up at the sky, our eyes see the blue color that has been scattered in all directions by the tiny molecules.\n4. The other colors of light, such as red and orange, continue to travel in a straight line through the atmosphere with less scattering, which is why they appear more vibrant when viewed from an overhead angle.\n\nThis effect is more pronounced during the daytime when the sun is high in the sky, and it's also more noticeable on clear days with minimal atmospheric interference. The color of the sky can change depending on various factors like atmospheric conditions, pollution, and time of day.\n\nWould you like to know more about Rayleigh scattering or the physics behind it?"

	c.Messages = append(c.Messages, m)
	return c, nil
}
