import { useEffect, useState } from "react";
import {
  ByAuthorContainer,
  Container,
  CooldownInputContainer,
  CreateSoundContainer,
  CreateSoundForm,
  InfoContainer,
  SoundTable,
  SoundTableActions,
  SoundTableContainer,
  SoundTableHelper,
  Title,
} from "../style/dashboard";
import { TitleDeploy } from "../util/TitleDeploy";
import { SoundMap, formatNumber, pageOf, notEmptyOrElse, Deployed } from "../util/shared";
import "react-toastify/dist/ReactToastify.css";
import { ToastContainer, toast } from "react-toastify";
import Axios, { AxiosResponse } from "axios";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlay, faTrash, faX } from "@fortawesome/free-solid-svg-icons";

const NUMBER_REGEX = /^\d+$/;

const DAY_UNIT = "day";
const HOUR_UNIT = "hour";
const MINUTE_UNIT = "minute";
const SECOND_UNIT = "second";
const MILLISECOND_UNIT = "millisecond";

const TranslateUnit = (unit: string, value: number): number => {
  switch (unit) {
    case DAY_UNIT:
      return value * 24 * 60 * 60 * 1000;
    case HOUR_UNIT:
      return value * 60 * 60 * 1000;
    case MINUTE_UNIT:
      return value * 60 * 1000;
    case SECOND_UNIT:
      return value * 1000;
    default:
      return value;
  }
};

const Upload = async (
  price: number,
  cooldown: number,
  name: string,
  formData: FormData
): Promise<boolean> => {
  return new Promise((resolve) => {
    Axios.post(`http://localhost:9999/sound`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      params: {
        price,
        cooldown,
        name,
      },
    })
      .catch(() => resolve(false))
      .then(() => resolve(true));
  });
};

const ToastError = (child: JSX.Element) => {
  toast(child, {
    style: {
      backgroundColor: "rgb(195 83 83)",
    },
  });
};

const BoldSuccessStyle = {
  fontWeight: "bold",
  color: "#2d8538",
};

const ToastSuccess = (child: JSX.Element) => {
  toast(child, {
    style: {
      backgroundColor: "#54b961",
    },
  });
};

function Actions({
  onPlay,
  onDelete,
}: {
  onPlay: () => void;
  onDelete: () => void;
}): JSX.Element {
  return (
    <>
      <button onClick={onPlay}>
        <FontAwesomeIcon icon={faPlay} />
      </button>
      <button onClick={onDelete}>
        <FontAwesomeIcon icon={faTrash} />
      </button>
    </>
  );
}

type NewAudioStructure = {
  price: string;
  cooldown: string;
  cooldownUnit: string;
  name: string;
  file: File | null;
};

// TODO: set up play action (http://localhost:9999/sound/test/:id - POST)
export default function Dashboard() {
  const [isCreating, setIsCreating]       = useState(false);
  const [isDeleting, setIsDeleting]       = useState(false);
  const [maxSoundsPage, setMaxSoundsPage] = useState(1);
  const [soundsPage, setSoundsPage]       = useState(1);
  const [sounds, setSounds]               = useState<SoundMap>({});
  const [newAudio, setNewAudio]           = useState<NewAudioStructure>(
    {
      price: "",
      cooldown: "",
      cooldownUnit: SECOND_UNIT,
      name: "",
      file: null
    }
  );

  useEffect(() => {
    Axios.get("http://localhost:9999/sounds")
      .catch(console.log)
      .then((it) => {
        const response = it as AxiosResponse<any, any>;
        setSounds((old) => {
          const sounds = { ...old, ...response.data };
          updateMaxSoundsPage(sounds);
          return sounds;
        });
      });
  }, []);

  const updateMaxSoundsPage = (it: SoundMap) => {
    const keys = Object.keys(it);
    setMaxSoundsPage(_ => {
      const value = keys.length == 0 ? 1 : Math.ceil(keys.length / 5);
      if (soundsPage > value) {
        changeSoundsPage(value);
      }
      return value;
    });
  };

  const changeSoundsPage = (page: number) => {
    if (soundsPage == page || page > maxSoundsPage) {
      return;
    }

    if (page <= 0) {
      if (soundsPage != 1) {
        setSoundsPage(1);
      }
      return;
    }

    setSoundsPage(page);
  };

  const deleteSound = (id: string) => {
    if (isDeleting) {
      return;
    }

    setIsDeleting(true);

    setSounds(old => {
      const newVal: SoundMap = {...old};
      delete newVal[id];
      updateMaxSoundsPage(newVal);
      return newVal;
    });

    Axios
      .delete(`http://localhost:9999/sound/${id}`)
      .catch(() => ToastError(<p>Failed deleting sound. Try refreshing the page.</p>))
      .then(res => {
        if (res !== undefined) {
          ToastSuccess(<p>Successfully deleted the sound <span style={BoldSuccessStyle}>{id}</span></p>);
        }
        setIsDeleting(false);
      });
  };

  const testSound = (id: string) => {
    Axios
      .post(`http://localhost:9999/sound/test/${id}`)
      .catch(() => ToastError(<p>Failed during deployment of test sound. Try refreshing the page.</p>))
      .then(res => {
        if (res !== undefined) {
          ToastSuccess(<p>Successfully deployed test sound of <span style={BoldSuccessStyle}>{id}</span>.</p>)
        }
      });
  };

  function updateNewAudio(mod: (struct: NewAudioStructure) => void): void {
    setNewAudio(old => {
      const newVal: NewAudioStructure = { ...old };
      mod(newVal);
      return newVal;
    });
  };

  return (
    <TitleDeploy title="Dashboard">
      <>
        <Container>
          <Title>Dashboard</Title>
          <SoundTableContainer>
            <SoundTable>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Price</th>
                  <th>Cooldown</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {notEmptyOrElse(pageOf(sounds, soundsPage, 5), () =>
                  Object.keys(sounds)
                ).map((key) => {
                  const sound = sounds[key];
                  if (!sound) {
                    return <></>;
                  }
                  return (
                    <tr key={sound.file_name}>
                      <td>{key}</td>
                      <td>{sound.price}</td>
                      <td>{formatNumber(sound.cooldown)}</td>
                      <td>
                        <SoundTableActions>
                          <Actions
                            onPlay={() => testSound(key)}
                            onDelete={() => deleteSound(key)}
                          />
                        </SoundTableActions>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </SoundTable>
            <SoundTableHelper>
              <label onClick={() => changeSoundsPage(soundsPage - 1)}>
                &#8249;
              </label>
              <button onClick={() => setIsCreating(true)}>Create</button>
              <label onClick={() => changeSoundsPage(soundsPage + 1)}>
                &#8250;
              </label>
            </SoundTableHelper>
          </SoundTableContainer>
        </Container>
        <CreateSoundContainer style={{display: isCreating ? "flex" : "none"}}>
          <CreateSoundForm buttonBackground={newAudio.file !== null ? "#5cc769" : "#eb5f5f"}>
            <header>
              <div>
                <FontAwesomeIcon icon={faX} onClick={() => setIsCreating(false)} />
              </div>
            </header>
            <div>
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

                updateNewAudio(it => it.file = file);
              }}/>
              <label htmlFor="choose-sound-to-upload">Click me to select file</label>
              <InfoContainer>
                <input type="text" placeholder="Name" onChange={change => updateNewAudio(it => it.name = change.target.value)} />
                <input type="text" placeholder="Price" onChange={change => updateNewAudio(it => it.price = change.target.value)} />
                <CooldownInputContainer>
                  <input type="text" placeholder="Cooldown" onChange={change => updateNewAudio(it => it.cooldown = change.target.value)} />
                  <select 
                    name="cooldown-units" 
                    id="cooldown-units" 
                    value={newAudio.cooldownUnit} 
                    onChange={change => updateNewAudio(it => it.cooldownUnit = change.target.value)}>
                    <option value={DAY_UNIT}>Day</option>
                    <option value={HOUR_UNIT}>Hour</option>
                    <option value={MINUTE_UNIT}>Minute</option>
                    <option value={SECOND_UNIT}>Second</option>
                    <option value={MILLISECOND_UNIT}>Millisecond</option>
                  </select>
                </CooldownInputContainer>
              </InfoContainer>
              <button disabled={newAudio.file === null} onClick={async () => {
                const selectedFile = newAudio.file;
                if (selectedFile === null) {
                  ToastError(<p>You must select an Audio file.</p>);
                  return;
                }

                const newAudioName = newAudio.name;
                if (newAudioName == "") {
                  ToastError(<p>Audio name must NOT be empty.</p>);
                  return
                }

                const newAudioPrice = newAudio.price;
                if (newAudioPrice == "" || !NUMBER_REGEX.test(newAudioPrice)) {
                  ToastError(<p>Audio Price is either invalid or missing - make sure it's a number with no decimals!</p>);
                  return
                }

                const newAudioCooldown = newAudio.cooldown;
                if (newAudioCooldown == "" || !NUMBER_REGEX.test(newAudioCooldown)) {
                  console.log(`'${newAudioCooldown}'`);
                  ToastError(<p>Audio Cooldown is either invalid or missing - make sure it's a number with no decimals!</p>);
                  return;
                }

                const formData = new FormData();
                formData.append("file", selectedFile);

                const audio: Deployed = {
                  price: parseInt(newAudioPrice),
                  file_name: selectedFile.name,
                  cooldown: TranslateUnit(newAudio.cooldownUnit, parseInt(newAudioCooldown)),
                  last_used: 0
                };

                const result = await Upload(
                  audio.price, 
                  audio.cooldown, 
                  newAudioName,
                  formData
                );

                if (result) {
                  setSounds((old) => {
                    const sounds = {
                      ...old,
                      [newAudioName]: audio,
                    };
                    updateMaxSoundsPage(sounds);
                    return sounds;
                  });
                  ToastSuccess(<p>You have added the Audio <span style={BoldSuccessStyle}>{newAudioName}</span> to the roster with a price of <span style={BoldSuccessStyle}>{newAudioPrice}</span>.</p>)
                } else {
                  ToastError(<p>Failed to upload the new Audio... Perhaps it already exists?</p>)
                }
              }}>{newAudio.file !== null ? `Add "${newAudio.file.name}" to the roster` : "None Selected"}</button>
            </div>
            <footer/>
          </CreateSoundForm>
        </CreateSoundContainer>
        <ToastContainer
          position="bottom-center"
          bodyStyle={{ color: "#fff" }}
          hideProgressBar={true}
          autoClose={3000}
        />
        <ByAuthorContainer>
          <p>Made with</p>
          <img src="https://pbs.twimg.com/media/FQd_mVSWQAQCW3U.jpg" />
          <p>by</p>
          <a href="https://twitter.com/oliwer_lindell" target="_blank" rel="noreferrer">Oliwer</a>
        </ByAuthorContainer>
      </>
    </TitleDeploy>
  );
}
