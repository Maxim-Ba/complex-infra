import { Component } from "@angular/core";
import { ChatComponent } from "../../widgets/chat/chat.component";

@Component({
  imports:[ChatComponent],
  template:"<app-chat/>"
})
export class TestChatComponent{

}
