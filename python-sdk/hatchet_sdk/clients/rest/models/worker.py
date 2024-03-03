# coding: utf-8

"""
    Hatchet API

    The Hatchet API

    The version of the OpenAPI document: 1.0.0
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


from __future__ import annotations
import pprint
import re  # noqa: F401
import json

from datetime import datetime
from pydantic import BaseModel, Field, StrictStr
from typing import Any, ClassVar, Dict, List, Optional
from hatchet_sdk.clients.rest.models.api_resource_meta import APIResourceMeta
from hatchet_sdk.clients.rest.models.step_run import StepRun
from typing import Optional, Set
from typing_extensions import Self

class Worker(BaseModel):
    """
    Worker
    """ # noqa: E501
    metadata: APIResourceMeta
    name: StrictStr = Field(description="The name of the worker.")
    last_heartbeat_at: Optional[datetime] = Field(default=None, description="The time this worker last sent a heartbeat.", alias="lastHeartbeatAt")
    actions: Optional[List[StrictStr]] = Field(default=None, description="The actions this worker can perform.")
    recent_step_runs: Optional[List[StepRun]] = Field(default=None, description="The recent step runs for this worker.", alias="recentStepRuns")
    __properties: ClassVar[List[str]] = ["metadata", "name", "lastHeartbeatAt", "actions", "recentStepRuns"]

    model_config = {
        "populate_by_name": True,
        "validate_assignment": True,
        "protected_namespaces": (),
    }


    def to_str(self) -> str:
        """Returns the string representation of the model using alias"""
        return pprint.pformat(self.model_dump(by_alias=True))

    def to_json(self) -> str:
        """Returns the JSON representation of the model using alias"""
        # TODO: pydantic v2: use .model_dump_json(by_alias=True, exclude_unset=True) instead
        return json.dumps(self.to_dict())

    @classmethod
    def from_json(cls, json_str: str) -> Optional[Self]:
        """Create an instance of Worker from a JSON string"""
        return cls.from_dict(json.loads(json_str))

    def to_dict(self) -> Dict[str, Any]:
        """Return the dictionary representation of the model using alias.

        This has the following differences from calling pydantic's
        `self.model_dump(by_alias=True)`:

        * `None` is only added to the output dict for nullable fields that
          were set at model initialization. Other fields with value `None`
          are ignored.
        """
        excluded_fields: Set[str] = set([
        ])

        _dict = self.model_dump(
            by_alias=True,
            exclude=excluded_fields,
            exclude_none=True,
        )
        # override the default output from pydantic by calling `to_dict()` of metadata
        if self.metadata:
            _dict['metadata'] = self.metadata.to_dict()
        # override the default output from pydantic by calling `to_dict()` of each item in recent_step_runs (list)
        _items = []
        if self.recent_step_runs:
            for _item in self.recent_step_runs:
                if _item:
                    _items.append(_item.to_dict())
            _dict['recentStepRuns'] = _items
        return _dict

    @classmethod
    def from_dict(cls, obj: Optional[Dict[str, Any]]) -> Optional[Self]:
        """Create an instance of Worker from a dict"""
        if obj is None:
            return None

        if not isinstance(obj, dict):
            return cls.model_validate(obj)

        _obj = cls.model_validate({
            "metadata": APIResourceMeta.from_dict(obj["metadata"]) if obj.get("metadata") is not None else None,
            "name": obj.get("name"),
            "lastHeartbeatAt": obj.get("lastHeartbeatAt"),
            "actions": obj.get("actions"),
            "recentStepRuns": [StepRun.from_dict(_item) for _item in obj["recentStepRuns"]] if obj.get("recentStepRuns") is not None else None
        })
        return _obj

