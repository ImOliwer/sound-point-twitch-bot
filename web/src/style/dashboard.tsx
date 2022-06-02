import styled from "styled-components";

type AdditionContainerProps = {
  buttonBackground: string
}

export const Container = styled.div`
  width: 100%;
  height: 100vh;
  background: #26262a;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
`;

export const Title = styled.h2`
  font-weight: bold;
  font-size: 36px;
  color: #f0f0f0;
  margin-bottom: 50px;
`;

export const SoundTableContainer = styled.div`
  display: flex;
  flex-direction: column;
  background: #17161C;
  padding: 20px;
  border-radius: 7.5px;
`;

export const SoundTable = styled.table`
  border-collapse: collapse;
  table-layout: fixed;

  & th {
    padding: 10px;
    text-align: start;
    text-transform: uppercase;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    font-size: 14px;
    color: #4CDA8D;
  }

  & th:last-child {
    text-align: center;
  }

  & td {
    padding: 10px;
    font-size: 15px;
    color: #ffffff8a;
  }

  & td:last-child button:not(:last-child) {
    margin-right: 10px;
  }

  & th, & td {
    width: 120px;
    max-width: 120px;
  }

  & td:not(:last-child) {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  & tbody tr:nth-child(odd) {
    background: rgba(255, 255, 255, 0.01);
  }
`;

export const SoundTableActions = styled.div`
  max-width: 120px;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;

  & button {
    outline: none;
    border: none;
    padding: 7px 10px 7px 10px;
    border-radius: 50%;
  }

  & button:first-child {
    background: #07bc0c;#f13939
  }

  & button:first-child:hover {
    background: #1acd1f;
  }

  & button:last-child {
    background: #f13939;
  }

  & button:last-child:hover {
    background: #f15252;
  }

  & button:first-child, & button:last-child {
    color: #fff;
    cursor: pointer;
  }

  & button:first-child,
  & button:last-child,
  & button:first-child:hover, 
  & button:last-child:hover {
    transition: ease-in-out .3s background;
  }
`;

export const SoundTableHelper = styled.div`
  width: 100%;
  padding-top: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;

  & label {
    padding: 2.5px 10px 2.5px 10px;
    background: #4e4e6421;
    color: #4CDA8D;
    font-size: 18px;
    border-radius: 50%;
    text-align: center;
    cursor: pointer;
    transition: ease-in-out .3s background;
    -webkit-touch-callout: none;
    -webkit-user-select: none;
     -khtml-user-select: none;
       -moz-user-select: none;
        -ms-user-select: none;
            user-select: none;
  }

  & label:hover {
    background: #4e4e6436;
    transition: ease-in-out .3s background;
  }

  & button {
    height: 100%;
    padding: 0 15px 0 15px;
    outline: none;
    border: none;
    border-radius: 5px;
    background: #4CDA8D;
    color: 2C2839;
    font-size: 12px;
    font-weight: bold;
    text-transform: uppercase;
    cursor: pointer;
    transition: ease-in-out .3s background;
  }

  & button:hover {
    background: #29df7c;
    transition: ease-in-out .3s background;
  }
`;

export const CreateSoundContainer = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  justify-content: center;
  align-items: center;
  background: rgba(0, 0, 0, 0.5);
`;

export const CreateSoundForm = styled.div<AdditionContainerProps>`
  width: 15%;
  height: 35%;
  padding: 15px;
  background: #e3e3e3;
  border-radius: 7.5px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-direction: column;

  & > header, & > footer {
    width: 100%;
    height: 10%;
  }

  & > header {
    display: flex;
    justify-content: end;
    align-items: center;

    & > div {
      padding: 7px;
      width: 16px;
      background: #c9c9c9;
      border-radius: 50%;
      display: flex;
      justify-content: center;
      align-items: center;
      transition: ease-in-out .3s background;
      cursor: pointer;

      &:hover {
        background: #a1a0a0;
        transition: ease-in-out .3s background;
      }

      & > * {
        width: 10px;
      }
    }
  }

  & > div {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;

    & > label {
      font-size: 16px;
      color: #3b3940;
      padding-bottom: 3.5px;
      cursor: pointer;
      border-bottom: 2px solid #3b3940;
      margin-bottom: 15px;
    }

    & > button {
      padding: 10px;
      width: 75%;
      border: none;
      outline: none;
      color: #fff;
      border-radius: 5px;
      text-transform: uppercase;
      font-weight: bold;
      background: ${props => props.buttonBackground || "#eb5f5f"};
      margin-top: 15px;
      cursor: pointer;
    }

    & > input {
      outline: none;
      border: none;
      padding: 10px;
      width: 50%;
    }
  }
`;

export const InfoContainer = styled.div`
  width: 75%;
  display: flex;
  flex-direction: column;

  & > input {
    margin-bottom: 5px;
  }
`;

export const CooldownInputContainer = styled.div`
  width: 100%;
  display: flex;

  & > input {
    width: 55%;
    margin-right: 5px;
  }

  & > select {
    flex: 1;
    border-radius: 5px;
    outline: none;
    border: none;
  }
`;

export const ByAuthorContainer = styled.div`
  position: absolute;
  top: 97.5%;
  left: 0.5%;  
  display: flex;
  color: rgba(255, 255, 255, 0.4);

  & img {
    width: 16px;
    height: 16px;
    text-align: center;
  }

  & a {
    text-decoration: none;
    color: #4CDA8D;
    font-weight: semi-bold;
  }

  & a:focus {
    color: #4CDA8D !important;
  }

  & p {
    margin-right: 5px;
  }

  & p:not(:first-child) {
    margin-left: 5px;
  }
`;