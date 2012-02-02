package LoginServer

import (
	. "SG"
	D "Data"
)

func OnWelcome(c *LClient, p *SGPacket) {
	c.Log().Println_Debug("OnWelcome Packet")
	packet := NewPacket2(15)
	packet.WriteHeader(CSM_WELCOME)
	packet.WriteUInt16(46381)
	packet.WriteUInt16(168)
	packet.WSkip(1)
	c.Send(packet)
}

func OnWelcome2(c *LClient, p *SGPacket) {
	c.Log().Printf_Debug("OnWelcome2 Packet")
	switch p.ReadByte() {
	case 0:
		packet := NewPacket2(20)
		packet.WriteHeader(SM_SENDIP)
		packet.Index--
		ip := []byte(Server.WANAddr.IP.To4()) 
		packet.WriteBytes([]byte{ip[3], ip[2], ip[1], ip[0]})
		packet.WriteUInt16(uint16(Server.WANAddr.Port))
		packet.WriteByte(0)
		packet.WriteByte(0)
		c.Send(packet)
		break
	case 1:
		packet := NewPacket2(200)
		packet.WriteHeader(SM_WELCOME2)
		packet.WriteBytes([]byte{0x00, 0x21, 0x78, 0x9C, 0x63, 0x64, 0x70, 0xD3, 0xAF, 0x71, 0xE0, 0xDE, 0xC1, 0x18, 0x9C, 0x91, 0x58, 0x52, 0x92, 0x5A, 0x94, 0x9A, 0xE2, 0x9E, 0x98, 0x93, 0x58, 0x51, 0xC9, 0x00, 0x00, 0x5D, 0x17, 0x08, 0x01})
		c.Send(packet)
		break
	}
}

func OnPlanetDataRequest(c *LClient, p *SGPacket) {
	c.Log().Printf_Debug("OnPlanetDataRequest Packet")
	p = NewPacket2(200)
	p.WriteHeader(CS_PLANET_DATA)

	p.WriteByte(1) // planet number loop

	p.WriteInt16(0)
	p.WriteString("Hell")
	p.WriteByte(2) //look
	p.WriteByte(7)
	p.WriteInt32(0)
	p.WSkip(3)               //angle 3 bytes (x,y,z) rotation 
	p.WriteInt16(100)        //radius
	p.WriteInt16(100)        //radius
	p.WriteInt32(3 * 100000) //location
	p.WriteInt16(10)
	p.WriteByte(1)
	p.WriteByte(1)     //available
	p.WriteInt32(0x00) //faction releated something
	p.WriteByte(0)     //faction disabled number loop
	//p.WriteInt32(0)

	p.WriteInt16(0) //default planet

	c.Send(p)
}

func OnRegister(c *LClient, p *SGPacket) {
	user := p.ReadString()
	pass := p.ReadString()
	email := p.ReadString()
	c.Log().Printf_Debug("OnRegister Packet User Register: User(%s) EMail(%s)", user, email)

	ec, s := D.CheckUser(user, email, pass)
	if ec == 1 {
		c.TempUser = &D.User{"", user, email, pass}
	}
	SendMessage(c, ec, s)
}

func OnFriendSelect(c *LClient, p *SGPacket) {
	user := p.ReadString()
	c.Log().Printf_Debug("OnFriendSelect Packet: User(%s)", user)

	//55 00 - not found
	//55 01 00 00 00 0C - found and faction

	packet := NewPacket2(20)
	packet.WriteHeader(SM_FRIEND_SELECT)
	packet.WriteByte(0)
	c.Send(packet)
}

func OnRegisterDone(c *LClient, p *SGPacket) {
	c.Log().Printf_Debug("OnRegisterDone: % #X\n", p.Buffer)
	if c.TempUser != nil {

		player := D.NewPlayer()
		player.Faction = p.ReadInt32()
		player.Avatar = p.ReadByte()
		p.RSkip(1) //avatar twice
		player.Tactics = p.ReadByte()
		player.Clout = p.ReadByte()
		player.Education = p.ReadByte()
		player.MechApt = p.ReadByte()
		player.Name = c.TempUser.User

		c.TempUser.ID = D.NewID()

		if D.RegisterUser(c.TempUser) == true {
			player.ID = D.NewID()
			player.UserID = c.TempUser.ID
			player.SetDefaultStats()
			D.RegisterPlayer(player)
		}
		c.TempUser = nil
	}
}

func OnLoginWelcome(c *LClient, p *SGPacket) {
	c.Log().Printf_Debug("OnLoginWelcome Packet")
	packet := NewPacket2(20)
	packet.WriteHeader(SM_REG_METHOD)
	packet.WriteByte(1) //1 advanced reg , else simple reg
	packet.WriteByte(0)
	c.Send(packet)
}

func OnFactionDataRequest(c *LClient, p *SGPacket) {
	c.Log().Printf_Debug("OnFactionDataRequest Packet")
	packet := NewPacket2(1000)
	packet.WriteHeader(CSM_FACTION_DATA)
	packet.WriteBytes([]byte{0x23, 
		0x00, 0x00, 0x00, 0x63, 0x08, 0x4B, 0x61, 0x6D, 0x61, 0x6E, 0x61, 0x73, 0x6F, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x0A, 0xAE, 0x67, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x62, 0x08, 0x4B, 0x61, 0x6D, 0x61, 0x6E, 0x61, 0x73, 0x69, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x0A, 0xB1, 0x22, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x61, 0x08, 0x4B, 0x61, 0x6D, 0x61, 0x6E, 0x61, 0x73, 0x65, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x0A, 0xB1, 0x1C, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x60, 0x08, 0x4B, 0x61, 0x6D, 0x61, 0x6E, 0x61, 0x73, 0x61, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x0A, 0xAE, 0x60, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x5F, 0x08, 0x43, 0x6C, 0x61, 0x64, 0x69, 0x63, 0x65, 0x6F, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x09, 0x27, 0xC7, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x5E, 0x08, 0x43, 0x6C, 0x61, 0x64, 0x69, 0x63, 0x65, 0x69, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x09, 0x2A, 0x82, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x5D, 0x08, 0x43, 0x6C, 0x61, 0x64, 0x69, 0x63, 0x65, 0x65, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x09, 0x2A, 0x7C, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x5C, 0x08, 0x43, 0x6C, 0x61, 0x64, 0x69, 0x63, 0x65, 0x61, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x09, 0x27, 0xC0, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x5B, 0x06, 0x47, 0x75, 0x6C, 0x62, 0x61, 0x6F, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x07, 0xA1, 0x27, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x5A, 0x06, 0x47, 0x75, 0x6C, 0x62, 0x61, 0x69, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x07, 0xA3, 0xE2, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x59, 0x06, 0x47, 0x75, 0x6C, 0x62, 0x61, 0x65, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x07, 0xA3, 0xDC, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x58, 0x06, 0x47, 0x75, 0x6C, 0x62, 0x61, 0x61, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x07, 0xA1, 0x20, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x17, 0x07, 0x41, 0x67, 0x61, 0x72, 0x74, 0x68, 0x61, 0x00, 0x00, 0x0F, 0x64, 0x0F, 0x00, 0x00, 0x01, 0x91, 0x08, 0x24, 0x1B, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x57, 0x07, 0x5A, 0x65, 0x6E, 0x64, 0x6F, 0x77, 0x6F, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x06, 0x1A, 0x87, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x16, 0x08, 0x50, 0x61, 0x63, 0x69, 0x66, 0x69, 0x63, 0x61, 0x00, 0x00, 0xC8, 0x0A, 0x0A, 0x00, 0x00, 0x00, 0x6A, 0x01, 0x7A, 0xA2, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x56, 0x07, 0x5A, 0x65, 0x6E, 0x64, 0x6F, 0x77, 0x69, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x06, 0x1D, 0x42, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x15, 0x08, 0x41, 0x74, 0x6C, 0x61, 0x6E, 0x74, 0x69, 0x73, 0x00, 0x00, 0x19, 0x2D, 0xC8, 0x00, 0x00, 0x02, 0xC2, 0x18, 0x71, 0xD1, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x55, 0x07, 0x5A, 0x65, 0x6E, 0x64, 0x6F, 0x77, 0x65, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x06, 0x1D, 0x3C, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x54, 0x07, 0x5A, 0x65, 0x6E, 0x64, 0x6F, 0x77, 0x61, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x06, 0x1A, 0x80, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x53, 0x09, 0x4D, 0x6F, 0x6E, 0x6F, 0x6C, 0x69, 0x74, 0x68, 0x6F, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x03, 0x0D, 0x47, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x52, 0x09, 0x4D, 0x6F, 0x6E, 0x6F, 0x6C, 0x69, 0x74, 0x68, 0x65, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x03, 0x0F, 0xFC, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x51, 0x09, 0x4D, 0x6F, 0x6E, 0x6F, 0x6C, 0x69, 0x74, 0x68, 0x69, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x03, 0x10, 0x02, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x50, 0x09, 0x4D, 0x6F, 0x6E, 0x6F, 0x6C, 0x69, 0x74, 0x68, 0x61, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x03, 0x0D, 0x40, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x4F, 0x04, 0x54, 0x68, 0x6F, 0x6F, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x0C, 0x35, 0x07, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x4E, 0x04, 0x54, 0x68, 0x6F, 0x69, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x0C, 0x37, 0xC2, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x4D, 0x04, 0x54, 0x68, 0x6F, 0x65, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x0C, 0x37, 0xBC, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x4C, 0x04, 0x54, 0x68, 0x6F, 0x61, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x0C, 0x35, 0x00, 0x02, 0x0A, 0xC9, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x08, 0x08, 0x44, 0x65, 0x66, 0x65, 0x6E, 0x64, 0x65, 0x72, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x01, 0x8A, 0x2C, 0xFE, 0x00, 0x00, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x07, 0x08, 0x41, 0x74, 0x74, 0x61, 0x63, 0x6B, 0x65, 0x72, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x01, 0x8A, 0x2C, 0xFD, 0x00, 0x00, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x06, 0x07, 0x4E, 0x65, 0x75, 0x74, 0x72, 0x61, 0x6C, 0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x01, 0x8A, 0x2C, 0xFF, 0x00, 0x00, 0x9C, 0x40,
		0x00, 0x00, 0x00, 0x05, 0x09, 0x4D, 0x65, 0x72, 0x63, 0x65, 0x6E, 0x61, 0x72, 0x79, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x01, 0x8A, 0x2C, 0xFC, 0x00, 0x00, 0x9C, 0x40,
		//Pompeii
		0x00, 0x00, 0x00, 0x0C, 0x07, 0x50, 0x6F, 0x6D, 0x70, 0x65, 0x69, 0x69, 0x00, 0x01, 0x00, 0x46, 0x96, 0x00, 0x00, 0x00, 0x00, 0x1B, 0x92, 0xD5, 0x9C, 0x40,
		//Helike	
		0x00, 0x00, 0x00, 0x0B, 0x06, 0x48, 0x65, 0x6C, 0x69, 0x6B, 0x65, 0x00, 0x01, 0x37, 0xb4, 0x46, 0x00, 0x00, 0x00, 0x00, 0x11, 0x71, 0x23, 0x9C, 0x40,
		//Alien
		0x00, 0x00, 0x00, 0x00, 0x05, 0x41, 0x6C, 0x69, 0x65, 0x6E, 0x00, 0x01, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x9C, 0x40,
		//Troy
		0x00, 0x00, 0x00, 0x0D, 0x04, 0x54, 0x72, 0x6F, 0x79, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x1D, 0x07, 0x9C, 0x40,

		0x00, 0x00,
		0x00, 0x00})

	/* 

		23 // number of factions

		loop {
		00 00 00 63 //id
		08 4B 61 6D 61 6E 61 73 6F //name
		00
		00 //available to choose
		FF FF FF //rbg color
		00 0A AE 67 //nplanet
		02 //icon
		0A C9
		9C 40 }

		00 00 //if > 0, reads more something with nations
		00 00 //if > 0, reads more

		nplanet = planet num * 100000
		//0x00, 0x01, 0x86, 0xa0, = 100000 = planet 1
	*/

	c.Send(packet)
}

func OnLogin(c *LClient, p *SGPacket) {
	user := p.ReadString()
	pass := p.ReadString()
	c.Log().Printf_Info("OnLogin Packet User Register: User(%s)", user)

	ec, s, id := D.Login(user, pass)
	SendMessage(c, ec, s)
	if ec == 0 {
		//D.LoginQueue.Add(c.IP, id)
		addClient(c.IP, id)
		SendToGameServer(c, user)
	}
}
