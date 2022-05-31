import { useState } from "react";
import { AdditionContainer, Container, CooldownInputContainer, InfoContainer, Title } from "../style/dashboard";
import { TitleDeploy } from "../util/TitleDeploy";
import "react-toastify/dist/ReactToastify.css";
import { ToastContainer, toast } from "react-toastify";
import Axios from "axios";

const NUMBER_REGEX = /^\d+$/;

const DAY_UNIT = "day";
const HOUR_UNIT = "hour";
const MINUTE_UNIT = "minute";
const SECOND_UNIT = "second";
const MILLISECOND_UNIT = "millisecond";

const TranslateUnit = (unit: string, value: number): number => {
  switch (unit) {
    case DAY_UNIT: return value * 24 * 60 * 60 * 1000;
    case HOUR_UNIT: return value * 60 * 60 * 1000;
    case MINUTE_UNIT: return value * 60 * 1000;
    case SECOND_UNIT: return value * 1000;
    default: return value;
  }
}

const Upload = async (price: number, cooldown: number, name: string, formData: FormData): Promise<boolean> => {
  return new Promise((resolve) => {
    Axios
      .post(`http://localhost:9998/sound/upload`, formData, {
        headers: {
          "Content-Type": "multipart/form-data"
        },
        params: {
          price,
          cooldown,
          name,
        }
      })
      .catch(() => resolve(false))
      .then(() => resolve(true));
  });
};

const ToastError = (child: JSX.Element) => {
  toast(child, {
    style: {
      backgroundColor: "rgb(195 83 83)",
    }
  })
};

const ToastSuccess = (child: JSX.Element) => {
  toast(child, {
    style: {
      backgroundColor: "#54b961",
    }
  })
};

export default function Dashboard() {
  const [newAudioPrice, setNewAudioPrice] = useState("");
  const [newAudioCooldown, setNewAudioCooldown] = useState("");
  const [newAudioCooldownUnit, setNewAudioCooldownUnit] = useState(SECOND_UNIT);
  const [newAudioName, setNewAudioName] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  return (
    <TitleDeploy title="Dashboard">
      <>
        <Container>
          <Title>Dashboard</Title>
          <AdditionContainer buttonBackground={selectedFile !== null ? "#5cc769" : "#eb5f5f"}>
            <input type="file" accept="Audio/mp3" id="choose-sound-to-upload" hidden onChange={element => {
              const files = element.target.files;
              if (!files || files.length == 0) {
                return;
              }

              // ensure it's an "mp3" audio file
              const file = files[0];
              if (!file.name.endsWith("mp3")) {
                ToastError(<p>The file type must be of <span style={{fontWeight: "bold"}}>MP3</span>.</p>);
                return;
              }

              setSelectedFile(file);
            }}/>
            <label htmlFor="choose-sound-to-upload">Upload a Sound</label>
            <InfoContainer>
              <input type="text" placeholder="Name" onChange={it => setNewAudioName(it.target.value)} />
              <input type="text" placeholder="Price" onChange={it => setNewAudioPrice(it.target.value)} />
              <CooldownInputContainer>
                <input type="text" placeholder="Cooldown" onChange={it => setNewAudioCooldown(it.target.value)} />
                <select 
                  name="cooldown-units" 
                  id="cooldown-units" 
                  value={newAudioCooldownUnit} 
                  onChange={it => setNewAudioCooldownUnit(it.target.value)}>
                  <option value={DAY_UNIT}>Day</option>
                  <option value={HOUR_UNIT}>Hour</option>
                  <option value={MINUTE_UNIT}>Minute</option>
                  <option value={SECOND_UNIT}>Second</option>
                  <option value={MILLISECOND_UNIT}>Millisecond</option>
                </select>
              </CooldownInputContainer>
            </InfoContainer>
            <button disabled={selectedFile === null} onClick={async () => {
              if (selectedFile === null) {
                ToastError(<p>You must select an Audio file.</p>);
                return;
              }

              if (newAudioName == "") {
                ToastError(<p>Audio name must NOT be empty.</p>);
                return
              }

              if (newAudioPrice == "" || !NUMBER_REGEX.test(newAudioPrice)) {
                ToastError(<p>Audio Price is either invalid or missing - make sure it's a number with no decimals!</p>);
                return
              }

              if (newAudioCooldown == "" || !NUMBER_REGEX.test(newAudioCooldown)) {
                console.log(`'${newAudioCooldown}'`);
                ToastError(<p>Audio Cooldown is either invalid or missing - make sure it's a number with no decimals!</p>);
                return;
              }

              const formData = new FormData();
              formData.append("file", selectedFile);

              const result = await Upload(
                parseInt(newAudioPrice), 
                TranslateUnit(newAudioCooldownUnit, parseInt(newAudioCooldown)), 
                newAudioName, 
                formData
              );

              if (result) {
                ToastSuccess(<p>You have added the Audio <span style={{fontWeight: "bold"}}>{newAudioName}</span> to the roster with a price of <span style={{fontWeight: "bold"}}>{newAudioPrice}</span>.</p>)
              } else {
                ToastError(<p>Failed to upload the new Audio... Perhaps it already exists?</p>)
              }
            }}>{selectedFile !== null ? `Add "${selectedFile.name}" to the roster` : "None Selected"}</button>
          </AdditionContainer>
        </Container>
        <ToastContainer 
          position="bottom-center" 
          bodyStyle={{color: "#fff"}} 
          hideProgressBar={true} 
          autoClose={3000}
        />
      </>
    </TitleDeploy>
  )
}