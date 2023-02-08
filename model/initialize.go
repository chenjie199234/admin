package model

const AdminProjectID = "0,1"

// common node,all projects have this
const UserControl = ",1"   //projectid+this=user control's nodeid
const RoleControl = ",2"   //projectid+this=role control's nodeid
const ConfigControl = ",3" //projectid+this=config control's nodeid

// only the admin project has this
const Proxy = ",4" //projectid+this=proxy's nodeid
